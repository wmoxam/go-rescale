package main

import (
    "fmt"
    "github.com/nfnt/resize"
    "image/jpeg"
    "net/http"
    "strconv"
    "log"
    "runtime"
)

func print404(response http.ResponseWriter, msg string) {
  http.Error(response, msg, http.StatusNotFound)
  log.Print(msg)
}

func handler(response http.ResponseWriter, request *http.Request) {
  width, err := strconv.ParseUint(request.URL.Query()["width"][0], 10, 0)
  if err != nil {
    print404(response, "Could not determine width")
    return
  }
  height, err := strconv.ParseUint(request.URL.Query()["height"][0], 10, 0)
  if err != nil {
    print404(response, "Could not determine height")
    return
  }

  imageUrl := request.URL.Query()["image_url"][0]

  resp, err := http.Get(imageUrl)

  if err != nil || resp.StatusCode != 200 {
    resp.Body.Close()
    print404(response, fmt.Sprintf("Could not fetch image %s", imageUrl))
  } else {
    img, err := jpeg.Decode(resp.Body)

    resp.Body.Close()

    if err != nil {
      print404(response, fmt.Sprintf("Could not decode image %s", imageUrl))
      return
    }

    image := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

    response.Header().Set("Content-Type", "image/jpeg")
    jpeg.Encode(response, image, nil)

    log.Printf("Resized %s to %dx%d", imageUrl, width, height)
  }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
