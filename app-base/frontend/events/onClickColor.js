const onClickColor = ({ callback }) => {
    const colorBarEl = document.querySelector('.color-bar');
    colorBarEl.addEventListener('click', ({target}) => {
       if (!target.matches('.color-box')) {
           return;
       }
       const color = target.dataset.color;

        if (typeof callback === 'function') {
            callback(color);
        }

    });
};

export default onClickColor;