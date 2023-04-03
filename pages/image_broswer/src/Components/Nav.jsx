export default function Nav(props) {
    const { handles, index, rootPath, containerSize, jpegOpts, navRef } = props;

    let curIdxNew = index.cur + 1;
    if (curIdxNew > index.max) {
        curIdxNew = index.max;
    }

    const handleJpegToggle = () => {
        jpegOpts.setToJpeg(!jpegOpts.toJpeg);
    }

    return (
        <nav className="navbar bg-primary text-white" ref={navRef}>
            <div className="flex-1 m-2 text-xl font-mono">Images: {rootPath}</div>

            <div className="navbar-end">
                <div className="btn btn-primary btn-sm m-y-3" onClick={handleJpegToggle}>{jpegOpts.toJpeg ? 'Jpeg' : 'Original'}</div>

                <div className="flex flex-col text-sm m-y-3">
                    <div>Max H: {containerSize.height}</div>
                    <div>Max W: {containerSize.width}</div>
                </div>

                <button className="btn btn-primary btn-sm" onClick={handles.prev} >Previous</button>
                <div className="flex flex-col items-center mx-2">
                    <div className="text-md font-mono">{curIdxNew}</div>
                    <div className="text-md font-mono">{index.max}</div>
                </div>

                <button className="btn btn-primary btn-sm" onClick={handles.next}>Next</button>
                <button className="btn btn-primary btn-sm" onClick={handles.delete} >Delete</button>
                <button className="btn btn-primary btn-sm" onClick={handles.refresh} >Refresh</button>
            </div>
        </nav>
    )
}
