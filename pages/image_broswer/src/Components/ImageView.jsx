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
                    wrapperClass="h-full-important"
                    contentClass="h-full-important"
                >
                    <div className="h-full-important">
                        <img
                            src={imageSrc}
                            className="h-full-important object-contain"
                            ref={imageRef}
                            />
                    </div>
                </TransformComponent>
            </TransformWrapper>
    );
}
