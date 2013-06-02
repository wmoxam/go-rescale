go-rescale
----------

Image rescaler written in Go.

	$ make && make build
	$ ./server -port 3000

	$ wget http://localhost:3000/?width=150&height=0&image_url=http://wmoxam.com/images/skyline.jpg

Note: specifying a height or width of zero will maintain the aspect ratio.
