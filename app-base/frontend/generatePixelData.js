const generatePixelData = ({pixelSize , width, height,} ={}) => {
    const colMax = width/pixelSize;
    const rowMax = height/pixelSize;

    const pixelData = [];
    for (let rowIndex = 0; rowIndex < rowMax; rowIndex++){
        const row =[];
        for (let colIndex = 0; colIndex < colMax; colIndex++){
            const color = 'white';
            row.push(color);
        }
        pixelData.push(row);
    }
    return pixelData
};

export default generatePixelData;