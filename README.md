# walldl-go
A wallpaper downloader from pexels written in golang

## usage
```
git clone https://github.com/ferretcode/walldl ~/walldl
mkdir -p ~/walldl/wallpapers
export PATH=$PATH:/home/$USER/walldl/bin/
export WALLDL_API_KEY=pexels api key
walldl -c "category" -n "1"
```

## arguments
- `walldl -c ""` The category name
- `walldl -n "1"` The number of wallpapers to download (up to 80)
