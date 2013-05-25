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
    "math/rand"
    "time"
    "runtime"
)

func randomString (l int ) string {
    bytes := make([]byte, l)
    for i:=0 ; i<l ; i++ {
        bytes[i] = byte(randInt(65,90))
    }
    return string(bytes)
}

func randInt(min int , max int) int {
    return min + rand.Intn(max-min)
}


func resizeImage(id string, outputImage chan image.Image, width uint, height uint, img image.Image) {
  log.Print(id+"Resize start")
  outputImage <- resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
  log.Print(id+"Resize end")
}

func handler(w http.ResponseWriter, r *http.Request) {
  id := randomString(10)
  log.Print(id+"Got Request")
  params := regexp.MustCompile("/").Split(r.URL.Path[1:], 3)
  width, err := strconv.ParseUint(params[0], 10, 0)
  height, err := strconv.ParseUint(params[1], 10, 0)
  imageUrl := "http://www.thestar.com/content/dam/thestar/uploads/2013/5/24/1369442175062.jpg.size.xlarge.promo.jpg"

  log.Print(id+"Fetching Image")
  resp, err := http.Get(imageUrl)
  log.Print(id+"Got Image")

  if err != nil || resp.StatusCode != 200 {
    resp.Body.Close()
    http.StatusText(404)
    fmt.Fprintf(w, "Could not fetch image %s", imageUrl)
    log.Printf(id+"Could not fetch image %s", imageUrl)
  } else {
    log.Print(id+"Begin Decode")
    img, err := jpeg.Decode(resp.Body)
    log.Print(id+"End Decode")
    resp.Body.Close()
    if err != nil {
      http.StatusText(404)
      fmt.Fprintf(w, "Could not decode image %s", imageUrl)
      log.Printf(id+"Could not decode image %s", imageUrl)
    }

    var outputImage chan image.Image = make(chan image.Image)

    go resizeImage(id, outputImage, uint(width), uint(height), img)

    image := <- outputImage

    w.Header().Set("Content-Type", "image/jpeg")
    log.Print(id+"Sending Response")
    jpeg.Encode(w, image, nil)
  }
}

func main() {
    rand.Seed( time.Now().UTC().UnixNano())
    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
