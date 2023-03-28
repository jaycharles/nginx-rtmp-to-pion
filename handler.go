// Assumes you have an http router listenining on 127.0.0.1:8081 with a path of /nginx-rtmp-handler pointed to this handler

func NginxRtmpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		return
	}

	// Nginx will send a multipart body with a properties we want. See nginx-rtmp git for details
	// Keep in mind that your local FFMPEG connections will result in requests to this handler, so consider adding something to the URL string to identify ffmpeg clients from publishing clients if your're doing authentication.

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid params", http.StatusBadRequest)
		return
	}

	event := r.Form.Get("call")

	if event == "connect" {
		// Not doing anything for now, but this is where you'd authenticate the request before RTMP publishing is accepted at nginx.
		// The stream name won't be known at this point, so if you need to validate the stream name do so in the publish event handling.
	}

	if event == "publish" {
		// Start ffmpeg. If you need to validate the rtmp stream name, do so here and return non 200 to reject.
		url := "rtmp://127.0.0.1/" + r.Form.Get("app") + "/" + r.Form.Get("name")
		ingress := &IngressFfmpeg{}
		ingress.SetHost(url)
		// You'll want to either store reference to the IngressFfmpeg or block on a goroutine here or something. this is just an example.
	}

	if event == "publish_done" {
		// Nothing for now, but a good place to clean up
	}
}
