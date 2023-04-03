
export default function ImageView(props) {
    const { imageSrc } = props;

    return (
        <div className="h-full" id="img-container">
            <img
                src={imageSrc}
                className="max-h-full object-contain"
            />
        </div>
    );
}
