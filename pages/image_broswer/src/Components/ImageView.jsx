import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";

import ImageLoader from "./ImageLoad";

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
                    <ImageLoader src={imageSrc} imgRef={imageRef} />
                </TransformComponent>
            </TransformWrapper>
    );
}
