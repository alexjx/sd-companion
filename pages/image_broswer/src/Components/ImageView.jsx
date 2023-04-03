import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";

export default function ImageView(props) {
    const { imageSrc, imageRef } = props;

    return (
            <TransformWrapper
                doubleClick={{
                    mode: "reset",
                }}
                >
                <TransformComponent
                    wrapperClass="h-full"
                    contentClass="h-full"
                >
                    <img
                        src={imageSrc}
                        className="h-full"
                        ref={imageRef}
                    />
                </TransformComponent>
            </TransformWrapper>
    );
}
