import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";

export default function ImageView(props) {
    const { imageSrc } = props;

    return (
        <TransformWrapper
            doubleClick={{
                mode: "reset",
            }}
        >
            <TransformComponent>
                <div className="h-full" id="img-container">
                    <img
                        src={imageSrc}
                        className="max-h-full object-contain"
                    />
                </div>
            </TransformComponent>
        </TransformWrapper>
    );
}
