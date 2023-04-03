import Details from "./Details";
import ImageView from "./ImageView";

export default function Broswer(props) {
    const { files, root, curIdx, innerRef, containerSize, toJpeg } = props;

    function contentUrl(file) {
        if (!toJpeg) {
            console.log('not to jpeg')
            return `${root}/files/${file.path}`
        }
        // path need to be url escaped
        const escapedPath = encodeURIComponent(file.path)
        return `${root}/api/encoded?path=${file.path}&height=${containerSize.height}`
    }

    return (
        <div className="flex item-center justify-center h-full max-w-7xl" >
            {
                curIdx < files.length && (
                    <div className="flex flex-row justify-center">
                        <div className="mx-3" ref={innerRef}>
                            <ImageView imageSrc={contentUrl(files[curIdx])} />
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
