package client

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3"
)

func EvenNumberRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + (i * 2)
	}
	return a
}

var IngressFfmpegPorts []int

type IngressFfmpeg struct {
	url           string
	cmd           *exec.Cmd
	audioPort     int
	videoPort     int
	audioListener *net.UDPConn
	videoListener *net.UDPConn
}

func (ingress *IngressFfmpeg) SetHost(url string) {
	ingress.url = url
	ingress.StartFF()
}

func (ingress *IngressFfmpeg) Destroy() {
	if ingress.cmd != nil {
		ingress.cmd.Process.Kill()
	}
	ingress.audioListener.Close()
	ingress.videoListener.Close()
	ingress.releasePorts()
}

func (ingress *IngressFfmpeg) createListener(port int) (*net.UDPConn, error) {
	l, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port})
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (ingress *IngressFfmpeg) getPorts() {
	if len(IngressFfmpegPorts) == 0 {
		IngressFfmpegPorts = EvenNumberRange(4000, 5000)
	}
	var a, v int
	a, IngressFfmpegPorts = IngressFfmpegPorts[0], IngressFfmpegPorts[1:]
	v, IngressFfmpegPorts = IngressFfmpegPorts[0], IngressFfmpegPorts[1:]
	ingress.audioPort = a
	ingress.videoPort = v
}

func (ingress *IngressFfmpeg) releasePorts() {
	if ingress.audioPort != 0 {
		IngressFfmpegPorts = append(IngressFfmpegPorts, ingress.audioPort)
		ingress.audioPort = 0
	}
	if ingress.videoPort != 0 {
		IngressFfmpegPorts = append(IngressFfmpegPorts, ingress.videoPort)
		ingress.videoPort = 0
	}
}

func (ingress *IngressFfmpeg) StartFF() {
	// Create listeners
	ingress.getPorts()
	audioListener, err := ingress.createListener(ingress.audioPort)
	if err != nil {
		//
	}
	ingress.audioListener = audioListener

	videoListener, err := ingress.createListener(ingress.videoPort)
	if err != nil {
		//
	}
	ingress.videoListener = videoListener

	go func() {
		wd, err := os.Getwd()
		if err != nil {
			//
		}
		args := []string{
			"-i", ingress.url,
			"-an",
			"-c:v", "copy",
			"-f", "rtp",
			"rtp://127.0.0.1:" + fmt.Sprint(ingress.videoPort) + "?pkt_size=1200",
			"-vn",
			"-c:a", "libopus",
			"-f", "rtp",
			"rtp:/127.0.0.1:" + fmt.Sprint(ingress.audioPort) + "?pkt_size=1200",
		}
		ingress.cmd = exec.Command("ffmpeg", args...)
		ingress.cmd.Dir = wd
		go func() {
			defer ingress.releasePorts()
			runerr := ingress.cmd.Run()
			if runerr != nil {
				// this would be a nonzero exit from ffmpeg.
			}

		}()
		go ingress.readRtp(&codecs.H264Packet{}, 90000, ingress.videoPort, "video", ingress.videoListener)
		ingress.readRtp(&codecs.OpusPacket{}, 48000, ingress.audioPort, "audio", ingress.audioListener)
	}()
}

func (ingress *IngressFfmpeg) readRtp(depacketizer rtp.Depacketizer, sampleRate uint32, port int, kind string, listener *net.UDPConn) {

	// Read RTP packets forever and send them to the egress
	for {
		inboundRTPPacket := make([]byte, 1500) // UDP MTU
		packet := &rtp.Packet{}
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			break
		}

		if err = packet.Unmarshal(inboundRTPPacket[:n]); err != nil {
			panic(err)
		}

		if kind == webrtc.RTPCodecTypeAudio.String() {
			// Do something with the audio packet, such as writing to a tracklocal
		}

		if kind == webrtc.RTPCodecTypeVideo.String() {
			// Do something with the the video packet, such as writing to a tracklocal
		}
	}

}
