const fisrtDrawPixel = ({canvasEl, ctx, pixelData, pixelSize, canvasWidth, canvasHeight}) => {
    canvasEl.width = canvasWidth;
    canvasEl.height = canvasHeight;
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, canvasWidth, canvasHeight);
}

export default fisrtDrawPixel;