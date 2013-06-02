package main

import (
    "flag"
    "fmt"
    "github.com/nfnt/resize"
    "image"
    "image/jpeg"
    "image/png"
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
    img, format, err := image.Decode(resp.Body)

    resp.Body.Close()

    if err != nil {
      print404(response, fmt.Sprintf("Could not decode image %s", imageUrl))
      return
    }

    image := resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)

    response.Header().Set("Content-Type", "image/png")

    log.Printf("Resized %s to %dx%d", imageUrl, width, height)
    
    switch format {
    case "jpeg":
      jpeg.Encode(response, image, nil)
    default:
      png.Encode(response, image)
    }
  }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    var port = flag.Int("port", 8080, "Port to listen on")
    flag.Parse()

    log.Printf("Starting iamge rescaler on port %d", *port)

    http.HandleFunc("/", handler)
    http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
