import Details from "./Details";
import ImageView from "./ImageView";

export default function Broswer(props) {
    const { files, root, curIdx, innerRef, containerSize } = props;

    function contentUrl(file) {
        // return `${root}/files/${file.path}`
        // path need to be url escaped
        const escapedPath = encodeURIComponent(file.path)
        return `${root}/api/encoded?path=${file.path}&height=${containerSize.height}`
    }

    return (
        <div className="flex item-center justify-center bg-gray-800 h-full" ref={innerRef}>
            {
                curIdx < files.length && (
                    <div className="flex flex-row justify-center mx-3">
                        <div className="mx-3">
                            <ImageView imageSrc={contentUrl(files[curIdx])} />
                        </div>
                        <div className="mx-3" >
                            <Details file={files[curIdx]} />
                        </div>
                    </div>
                )
            }
        </div>
    )
}
