import { useEffect, useState } from "react";

import Details from "./Details";
import ImageView from "./ImageView";

export default function Broswer(props) {
    const { files, root, curIdx, innerRef, containerSize, toJpeg } = props;

    const [imgSrc, setImgSrc] = useState(null);

    function contentUrl(file) {
        let prefix = ''
        if (root) {
            prefix = `${root}`
        }

        if (!toJpeg) {
            return `${prefix}/files/${file.path}`
        }
        // path need to be url escaped
        const escapedPath = encodeURIComponent(file.path)
        return `${prefix}/api/encoded?path=${file.path}&height=${containerSize.height}`
    }

    useEffect(() => {
        if (files && files.length > 0 && curIdx < files.length) {
            setImgSrc(contentUrl(files[curIdx]));
        }
    }, [curIdx, files, toJpeg]);

    return (
        <div className="flex item-center justify-center h-full max-w-7xl" >
            {
                curIdx < files.length && (
                    <div className="flex flex-row justify-center">
                        <div className="mx-3" ref={innerRef}>
                            {imgSrc && <ImageView imageSrc={imgSrc} />}
                        </div>
                        <div className="mx-3 max-w-xs overflow-y-auto" >
                            <Details file={files[curIdx]} />
                        </div>
                    </div>
                )
            }
        </div>
    )
}
