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
    "strings"
    "log"
    "runtime"
)

func ignoreRequestHandler(response http.ResponseWriter, request *http.Request) {
  logRequest(request)
  print404(response, "Not found")
}

func logRequest(request *http.Request) {
  uri := request.URL.RequestURI()
  log.Printf("Handling Request %s", uri)
}

func print404(response http.ResponseWriter, msg string) {
  http.Error(response, msg, http.StatusNotFound)
  log.Print(msg)
}

func handler(response http.ResponseWriter, request *http.Request) {
  logRequest(request)
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
  algorithm_name := request.URL.Query()["algorythm"][0]

  imageUrl := request.URL.Query()["image_url"][0]

  resp, err := http.Get(imageUrl)

  if err != nil || resp.StatusCode != 200 {
    print404(response, fmt.Sprintf("Could not fetch image %s", imageUrl))
  } else {
    img, format, err := image.Decode(resp.Body)

    resp.Body.Close()

    if err != nil {
      print404(response, fmt.Sprintf("Could not decode image %s", imageUrl))
      return
    }

    algorithm := algorithmFromName(algorithm_name)

    image := resize.Resize(uint(width), uint(height), img, algorithm)

    log.Printf("Resized %s to %dx%d using %s", imageUrl, width, height, algorithm_name)

    switch format {
    case "jpeg":
      response.Header().Set("Content-Type", "image/jpeg")
      jpeg.Encode(response, image, nil)
    default:
      response.Header().Set("Content-Type", "image/png")
      png.Encode(response, image)
    }
  }
}

func algorithmFromName(name string) resize.InterpolationFunction {
  switch strings.ToLower(name) {
  case "nearest_neighbour":
    return resize.NearestNeighbor
  case "bilinear":
    return resize.Bilinear
  case "bicubic":
    return resize.Bicubic
  case "mitchell_netravali":
    return resize.MitchellNetravali
  case "lanczos2":
    return resize.Lanczos2
  case "lanczos3":
    return resize.Lanczos3
  default:
    return resize.NearestNeighbor
  }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    var port = flag.Int("port", 7000, "Port to listen on")
    flag.Parse()

    log.Printf("Starting image rescaler on port %d", *port)

    http.HandleFunc("/", handler)
    http.HandleFunc("/favicon.ico", ignoreRequestHandler)
    http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
