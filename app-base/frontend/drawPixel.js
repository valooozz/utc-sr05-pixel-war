/*const drawPixel = ({
    ctx,
    pixelData,
    pixelSize,
    canvasWidth,
    canvasHeight,
}) => {
    const colMax = canvasWidth/pixelSize;
    const rowMax = canvasHeight/pixelSize;

    for (let rowIndex = 0; rowIndex < rowMax; rowIndex++){
        for (let colIndex = 0; colIndex < colMax; colIndex++){
            if (pixelData[colIndex][rowIndex] != null) {
                ctx.fillStyle = pixelData[colIndex][rowIndex];
                ctx.fillRect(colIndex*pixelSize, rowIndex*pixelSize, pixelSize, pixelSize);
            }


        }

    }
}*/
const drawPixel = ({ctx, color, rowIndex, colIndex, pixelData, pixelSize}) => {
    console.log("x", rowIndex);
    pixelData[rowIndex][colIndex] = color;
    console.log("color", color)
    ctx.fillStyle = pixelData[rowIndex][colIndex];
    ctx.fillRect(rowIndex*pixelSize, colIndex*pixelSize, pixelSize, pixelSize);
};

export default drawPixel;