import { DateTime } from 'luxon';

export default function Nav(props) {
    const { handles, index, rootPath, containerSize } = props;

    let curIdxNew = index.cur + 1;
    if (curIdxNew > index.max) {
        curIdxNew = index.max;
    }

    return (
        <nav className="navbar bg-primary text-white">
            <div className="flex-1 m-2 text-xl font-mono">Images: {rootPath}</div>
            <div className="navbar-end">
                <div className="flex flex-col text-sm m-y-3">
                    <div>H: {containerSize.height}</div>
                    <div>W: {containerSize.width}</div>
                </div>
                <button className="btn btn-primary" onClick={handles.prev} >Previous</button>
                <div className="m-2 text-xl font-mono">{curIdxNew} / {index.max}</div>
                <button className="btn btn-primary" onClick={handles.next}>Next</button>
                <button className="btn btn-primary" onClick={handles.delete} >Delete</button>
                <button className="btn btn-primary" onClick={handles.refresh} >Refresh</button>
            </div>
        </nav>
    )
}
