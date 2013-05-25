package main

import (
    "fmt"
    "github.com/nfnt/resize"
    "image"
    "image/jpeg"
    "net/http"
    "regexp"
    "strconv"
    "log"
    "runtime"
)

func resizeImage(outputImage chan image.Image, width uint, height uint, img image.Image) {
  log.Print("Resize start")
  outputImage <- resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
  log.Print("Resize end")
}

func handler(w http.ResponseWriter, r *http.Request) {
  log.Print("Got Request")
  params := regexp.MustCompile("/").Split(r.URL.Path[1:], 3)
  width, err := strconv.ParseUint(params[0], 10, 0)
  height, err := strconv.ParseUint(params[1], 10, 0)
  imageUrl := "http://www.thestar.com/content/dam/thestar/uploads/2013/5/24/1369442175062.jpg.size.xlarge.promo.jpg"

  log.Print("Fetching Image")
  resp, err := http.Get(imageUrl)
  log.Print("Got Image")

  if err != nil || resp.StatusCode != 200 {
    resp.Body.Close()
    http.StatusText(404)
    fmt.Fprintf(w, "Could not fetch image %s", imageUrl)
    log.Printf("Could not fetch image %s", imageUrl)
  } else {
    log.Print("Begin Decode")
    img, err := jpeg.Decode(resp.Body)
    log.Print("End Decode")
    resp.Body.Close()
    if err != nil {
      http.StatusText(404)
      fmt.Fprintf(w, "Could not decode image %s", imageUrl)
      log.Printf("Could not decode image %s", imageUrl)
    }

    var outputImage chan image.Image = make(chan image.Image)

    go resizeImage(outputImage, uint(width), uint(height), img)

    image := <- outputImage

    w.Header().Set("Content-Type", "image/jpeg")
    log.Print("Sending Response")
    jpeg.Encode(w, image, nil)
  }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
