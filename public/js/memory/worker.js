function customImageTexture (imgdata, string) {
    const vcanvas = new OffscreenCanvas(250, 250)
    const vctx = vcanvas.getContext('2d')

    vctx.font = "180px sans-serif";
    vctx.textAlign = "center";
    vctx.textBaseline = 'middle';
    vctx.fillText(string, vcanvas.width/2, vcanvas.height/2);

    const cropped = cropImageFromCanvas(vctx)

    const canvas = new OffscreenCanvas(512, 512)
    const ctx = canvas.getContext("2d")

    const baseImage = new ImageData(imgdata, 512, 512);

    ctx.putImageData(baseImage, 0, 0);

    const placeholder = { x: 107, y: 2, width: 253, height: 251 }
    const left = parseInt( (placeholder.width - cropped.width)/2 + placeholder.x )
    const top = parseInt( (placeholder.height - cropped.height)/2 + placeholder.y )

    ctx.drawImage(cropped, left, top)

    canvas.convertToBlob({
        type: "image/webp",
        quality: 1
    }).then(blob => {
        const reader = new FileReader();
        reader.onload = event => {
            const dataurl = event.target.result;
            console.log("MB: " + byteCount(dataurl) / 1e+6);
            self.postMessage({id: string, base64: dataurl})
        }
        reader.readAsDataURL(blob)
    })
}

function byteCount(s) {
    return encodeURI(s).split(/%..|./).length - 1;
}

function cropImageFromCanvas(ctx) {
    var canvas = ctx.canvas, 
      w = canvas.width, h = canvas.height,
      pix = {x:[], y:[]},
      imageData = ctx.getImageData(0,0,canvas.width,canvas.height),
      x, y, index;
    
    for (y = 0; y < h; y++) {
      for (x = 0; x < w; x++) {
        index = (y * w + x) * 4;
        if (imageData.data[index+3] > 0) {
          pix.x.push(x);
          pix.y.push(y);
        } 
      }
    }
    pix.x.sort(function(a,b){return a-b});
    pix.y.sort(function(a,b){return a-b});
    var n = pix.x.length-1;
    
    w = 1 + pix.x[n] - pix.x[0];
    h = 1 + pix.y[n] - pix.y[0];
    var cut = ctx.getImageData(pix.x[0], pix.y[0], w, h);
  
    canvas.width = w;
    canvas.height = h;
    ctx.putImageData(cut, 0, 0);
          
    return canvas
}

self.addEventListener('message', event => {
    customImageTexture(event.data.baseImage, event.data.str);
});