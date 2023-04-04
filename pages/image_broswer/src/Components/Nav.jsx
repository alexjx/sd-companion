import { BsArrowBarLeft, BsArrowBarRight } from 'react-icons/bs';
import { MdOutlineAutoDelete } from 'react-icons/md';
import { IoRefreshCircleOutline } from 'react-icons/io5';
import { useEffect, useState } from 'react';

export default function Nav(props) {
    const { handles, index, rootPath, containerSize, jpegOpts, navRef } = props;

    const handleJpegToggle = () => {
        jpegOpts.setToJpeg(!jpegOpts.toJpeg);
    }

    const Pagination = () => {
        const [curIdx, setCurIdx] = useState(0);

        useEffect(() => {
            let curIdxNew = index.cur + 1;
            if (curIdxNew > index.max) {
                curIdxNew = index.max;
            }
            setCurIdx(curIdxNew);
        }, [index.cur, index.max]);

        useEffect(() => {
            const timer = setTimeout(() => {
                const i = curIdx - 1;
                if (i >= 0 && i < index.max) {
                    handles.setCurIdx(i);
                }
            }, 1000);
            return () => clearTimeout(timer);
        }, [curIdx]);

        const handleChange = (e) => {
            setCurIdx(e.target.value);
        }

        return (
            <>
                <button className="btn btn-primary btn-sm text-2xl font-black" onClick={handles.prev} >
                    <BsArrowBarLeft />
                </button>

                <div className="flex flex-col items-center mx-2">
                    <input type="text" value={curIdx} className="input input-ghost input-xs text-center text-white input-primary w-10 font-mono" onChange={handleChange} />
                    <input type="text" value={index.max} className="input input-ghost input-xs text-center text-white input-primary w-10 font-mono disabled" />
                </div>

                <button className="btn btn-primary btn-sm text-2xl font-black" onClick={handles.next}>
                    <BsArrowBarRight />
                </button>
            </>
        )
    };

    return (
        <nav className="navbar bg-primary text-white" ref={navRef}>
            <div className="flex-1 m-2 text-xl font-mono">Images: {rootPath}</div>

            <div className="navbar-end">
                <div className="btn btn-primary btn-sm m-y-3" onClick={handleJpegToggle}>{jpegOpts.toJpeg ? 'JPEG' : 'ORIG'}</div>

                <div className="flex flex-col text-sm m-y-3">
                    <div>Max H: {containerSize.height}</div>
                    <div>Max W: {containerSize.width}</div>
                </div>

                {Pagination()}

                <button className="btn btn-primary btn-sm text-2xl font-black" onClick={handles.delete} >
                    <MdOutlineAutoDelete />
                </button>
                <button className="btn btn-primary btn-sm text-2xl font-black" onClick={handles.refresh} >
                    <IoRefreshCircleOutline />
                </button>
            </div>
        </nav>
    )
}
