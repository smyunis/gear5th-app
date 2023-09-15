package adsinteractors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdView struct {
	ViewIDToken string
	AdPiece     adpiece.AdPiece
}

type AdsPool struct {
	cacheStore         application.KeyValueStore
	campaignRepository campaign.CampaignRepository
	adPieceRepository  adpiece.AdPieceRepository
	digitalSignService application.DigitalSignatureService
	logger             application.Logger
}

func NewAdsPool(
	cacheStore application.KeyValueStore,
	campaignRepository campaign.CampaignRepository,
	adPieceRepository adpiece.AdPieceRepository,
	digitalSignService application.DigitalSignatureService,
	logger application.Logger) AdsPool {
	return AdsPool{
		cacheStore,
		campaignRepository,
		adPieceRepository,
		digitalSignService,
		logger,
	}
}

func (p *AdsPool) Next(slot adslot.AdSlotType) (*AdView, error) {

	curindex, err := p.cacheStore.Get(adsPoolSlotIndexCacheKey(slot))
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			loadIndexErr := p.loadAdPieceIndexes()
			if loadIndexErr != nil {
				p.logger.Error("adspool/cachestore/load-adpiece-indexes", loadIndexErr)
				return &AdView{}, loadIndexErr
			}
		} else {
			p.logger.Error("adspool/cachestore/get-cur-index", err)
			return &AdView{}, err
		}
	}
	cur, _ := strconv.Atoi(curindex)

	ad, err := p.cacheStore.Get(adsPoolAdPieceCacheKey(slot, cur))
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			loadAdPiecesErr := p.loadAdPieces()
			if loadAdPiecesErr != nil {
				p.logger.Error("adspool/cachestore/load-adpiece", err)
				return &AdView{}, err
			}
			ad, err = p.cacheStore.Get(adsPoolAdPieceCacheKey(slot, cur))
			if err != nil {
				p.logger.Error("adspool/cachestore/get-next-ad-post-reload", err)
				return &AdView{}, err
			}
		} else {
			p.logger.Error("adspool/cachestore/get-next-ad", err)
			return &AdView{}, err
		}
	}

	nextAdPiece := &adpiece.AdPiece{}
	err = json.Unmarshal([]byte(ad), nextAdPiece)
	if err != nil {
		p.logger.Error("adspool/unmarsharl-json/next-adpiece", err)
		return &AdView{}, err
	}

	p.moveCursorIndex(slot, cur)

	viewID := shared.NewID()
	viewIDHash, err := p.digitalSignService.Generate(viewID.String())
	if err != nil {
		p.logger.Error("adspool/generate-view-id-hash", err)
	}

	err = p.cacheStore.Save(ViewIDCacheKey(viewID.String()), "1", 1*time.Hour)
	if err != nil {
		p.logger.Error("adspool/cachestore/save-generated-view-id", err)
	}

	return &AdView{
		AdPiece:     *nextAdPiece,
		ViewIDToken: viewIDHash,
	}, nil

}

func (p *AdsPool) moveCursorIndex(slot adslot.AdSlotType, cur int) {
	ps, err := p.cacheStore.Get(adsPoolLengthCacheKey(slot))
	if err != nil {
		ps = "1000"
	}
	poolLen, _ := strconv.Atoi(ps)

	if cur == poolLen-1 {
		cur = -1
	}

	err = p.cacheStore.Save(adsPoolSlotIndexCacheKey(slot), strconv.Itoa(cur+1), 0)
	if err != nil {
		p.logger.Error("adspool/cachestore/save-cur-index", err)
	}
}

func (p *AdsPool) loadAdPieces() error {

	err := p.cacheStore.Save(adsPoolLengthCacheKey(adslot.Box), "0", 0)
	err = p.cacheStore.Save(adsPoolLengthCacheKey(adslot.Vertical), "0", 0)
	err = p.cacheStore.Save(adsPoolLengthCacheKey(adslot.Horizontal), "0", 0)
	if err != nil {
		p.logger.Error("adspool/cachestore/save-adspool-length", err)
	}

	runningCampaigns, err := p.campaignRepository.RunningCampaigns()
	if err != nil {
		return err
	}

	for _, c := range runningCampaigns {
		adPieces, err := p.adPieceRepository.ActiveAdPiecesForCampaign(c.ID)
		if err != nil {
			if errors.Is(err, application.ErrEntityNotFound) {
				continue
			}
			p.logger.Error("adspool/get-adpieces-for-campaign", err)
			return err
		}

		for _, a := range adPieces {
			j, err := json.Marshal(a)
			if err != nil {
				p.logger.Error("adspool/marshal-adpiece-json", err)
				continue
			}

			slotLengthStr, _ := p.cacheStore.Get(adsPoolLengthCacheKey(a.SlotType))
			slotLength, _ := strconv.Atoi(slotLengthStr)

			err = p.cacheStore.Save(adsPoolAdPieceCacheKey(a.SlotType, slotLength), string(j), 24*time.Hour)
			if err != nil {
				p.logger.Error("adspool/save-adpiece-to-pool", err)
				continue
			}
			err = p.cacheStore.Save(adsPoolLengthCacheKey(a.SlotType), strconv.Itoa(slotLength+1), 0)
		}
	}

	return nil
}

func (p *AdsPool) loadAdPieceIndexes() error {
	err := p.cacheStore.Save(adsPoolSlotIndexCacheKey(adslot.Box), "0", 0)
	err = p.cacheStore.Save(adsPoolSlotIndexCacheKey(adslot.Vertical), "0", 0)
	err = p.cacheStore.Save(adsPoolSlotIndexCacheKey(adslot.Horizontal), "0", 0)
	if err != nil {
		p.logger.Error("adspool/index/cache-save", err)
		return err
	}
	return nil
}

func adsPoolAdPieceCacheKey(slot adslot.AdSlotType, index int) string {
	return fmt.Sprintf("adspool:%d:%d", slot, index)
}

func adsPoolSlotIndexCacheKey(slot adslot.AdSlotType) string {
	return fmt.Sprintf("adspool:%d:index", slot)
}

func adsPoolLengthCacheKey(slot adslot.AdSlotType) string {
	return fmt.Sprintf("adspool:%d:length", slot)
}

func ViewIDCacheKey(viewID string) string {
	return fmt.Sprintf("adview:%s", viewID)
}
