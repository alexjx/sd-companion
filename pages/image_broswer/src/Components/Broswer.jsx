import { useEffect, useState } from "react";

import Details from "./Details";
import ImageView from "./ImageView";

export default function Broswer(props) {
    const {
        files,
        root,
        curIdx,
        innerRef,
        toJpeg,
        imageRef,
        fromTrash,
        containerSize,
    } = props;

    const [imgSrc, setImgSrc] = useState(null);

    function contentUrl(file) {
        let prefix = ''
        if (root) {
            prefix = `${root}`
        }

        if (!toJpeg) {
            if (!fromTrash) {
                return `${prefix}/files/${file.path}`
            }
            return `${prefix}/trash/${file.path}`
        }

        // path need to be url escaped
        const escapedPath = encodeURIComponent(file.path)
        if (!fromTrash) {
            return `${prefix}/api/encoded?path=${escapedPath}&width=${containerSize.width}&height=${containerSize.height}`
        }
        return `${prefix}/api/encoded?path=${escapedPath}&trash=1&&width=${containerSize.width}&height=${containerSize.height}`
    }

    useEffect(() => {
        if (files && files.length > 0 && curIdx < files.length) {
            setImgSrc(contentUrl(files[curIdx]));
        }
    }, [curIdx, files, toJpeg]);

    return (
        <div className="flex flex-row justify-center container-h" ref={innerRef}>
            <div className="mx-3">
                {imgSrc && <ImageView imageSrc={imgSrc} imageRef={imageRef} />}
            </div>
            <div className="mx-3 max-w-sm overflow-x-hidden overflow-y-auto" >
                <Details file={files ? files[curIdx] : null} fromTrash={fromTrash} containerSize={containerSize} />
            </div>
        </div>
    )
}
