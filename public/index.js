document.addEventListener('DOMContentLoaded', () => {
    const modal = document.querySelector('#myModal');
    const playButton = document.querySelector('#playButton');
    const videoContainer = document.querySelector('#videoContainer');
    const videoElement = document.querySelector('#video');

    playButton.addEventListener('click', async () => {
        modal.style.display = 'none';
        videoContainer.removeAttribute('hidden');
        videoContainer.style.display = 'block';
        videoElement.crossOrigin = 'anonymous';
        videoElement.setAttribute("autoplay", 'true');
        videoElement.setAttribute("playsinline", 'true');
        videoElement.setAttribute("muted", 'true');
        const hlsUrl = '/output.m3u8';

        if (Hls?.isSupported?.()) {
            const hls = new Hls({'debug': true});
            hls.loadSource(hlsUrl);
            hls.attachMedia(videoElement);
            hls.on(Hls.Events.MANIFEST_PARSED, async () => {
                try {
                    await videoElement.play();
                } catch (error) {
                    console.error('Videoの再生に失敗しました:', error);
                }
            });
        } else if (videoElement.canPlayType?.('application/vnd.apple.mpegurl')) {
            videoElement.src = hlsUrl;
            videoElement.addEventListener('loadedmetadata', async () => {
                try {
                    await videoElement.play();
                } catch (error) {
                    console.error('Videoの再生に失敗しました:', error);
                }
            });
        } else {
            alert('お使いのブラウザはHLSをサポートしていません。');
        }
    });
});
