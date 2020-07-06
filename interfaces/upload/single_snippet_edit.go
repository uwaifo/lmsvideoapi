package fileupload

import (
	"fmt"
	"time"

	"github.com/jtguibas/cinema"
)

//LocalEdit

func SingleSnippetEdit(filename, snippetStart, snippetEnd string) (string, error) {
	//Test Cinema begin
	video, err := cinema.Load("temp_upload/" + filename)
	check(err)

	fmt.Println(snippetStart + ":" + snippetEnd)

	//Convert string arguments to time
	start, err := ParseTime(snippetStart)
	if err != nil {
		fmt.Println(err)
	}

	end, err := ParseTime(snippetEnd)
	if err != nil {
		fmt.Println(err)
	}

	video.Trim(10*time.Second, 20*time.Second) // trim video from 10 to 20 seconds
	video.SetStart(start)                      // trim first second of the video
	video.SetEnd(end)                          // keep only up to 9 seconds
	//video.CommandLine()
	video.SetSize(400, 300)            // resize video to 400x300
	video.Crop(0, 0, 200, 200)         // crop rectangle top-left (0,0) with size 200x200
	video.SetSize(400, 400)            // resize cropped 200x200 video to a 400x400
	video.SetFPS(48)                   // set the output framerate to 48 frames per second
	video.SetBitrate(200_000)          // set the output bitrate of 200 kbps
	video.Render("edited_" + filename) // note format conversion by file extension

	//video.CommandLine()
	// you can also generate the command line instead of applying it directly
	fmt.Println("FFMPEG Command", video.CommandLine("test_output.mov"))
	return "editedAsset", err

	//Test cinema ends

}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
