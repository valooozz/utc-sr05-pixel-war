const onClickPixel= ({ ctx, color, canvasEl, pixelSize, pixelData, callback}) => {
    canvasEl.addEventListener('click', (evt) =>{
        console.log("evt", evt);
        const x = evt.offsetX;
        const y = evt.offsetY;

        const rowIndex = Math.floor(x/pixelSize);
        const colIndex = Math.floor(y/pixelSize);

        if (typeof callback === 'function') {
            callback(ctx, rowIndex, colIndex, pixelSize);
        }

    });
};

export default onClickPixel;
