# PWD must be project root
# needs tailwindcss executable https://tailwindcss.com/blog/standalone-cli
tailwindcss -i ./web/tailwind.init.css -o ./web/public/styles/tailwind.dist.css -c ./web/tailwind.config.js $1