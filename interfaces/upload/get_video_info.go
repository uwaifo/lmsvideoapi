package fileupload

/*
import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// getVideoInfo get audio and video duration in secs with video dimension as well
func GetVideoInfo(videoPath string) (*MediaInfo, error) {
	info := &MediaInfo{0, 0, 0, "", ""}
	// ffprobe -v quiet -print_format json -show_streams -show_format
	// ffprobe -v quiet -show_entries stream=width,height -of default=noprint_wrappers=1:nokey=1

	// get dimension
	if dimensionOutput, err := exec.Command(
		commands.FFMPEG.FFProbe,
		"-v", "quiet",
		"-show_entries", "stream=width,height",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("exec ffprobe %s with err: %s", videoPath, err)
	} else if info.Width, info.Height, err = getWidthAndHeightFromBytes(dimensionOutput); err != nil {
		return nil, fmt.Errorf("failed get width & height from %s with err %s", dimensionOutput, err)
	}

	// get duration
	if d, err := getDuration(videoPath); err == nil {
		info.Duration = d
	} else {
		return nil, err
	}
	return info, nil
}


*/
