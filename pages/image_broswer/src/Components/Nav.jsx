import { BsArrowBarLeft, BsArrowBarRight } from 'react-icons/bs';
import { MdOutlineAutoDelete } from 'react-icons/md';
import { IoRefreshCircleOutline } from 'react-icons/io5';
import { useEffect, useState, useRef } from 'react';
import axios from 'axios';

export default function Nav(props) {
    const {
        handles,
        index,
        currentDir,
        containerSize,
        jpegOpts,
        navRef,
    } = props;

    const curIdxRef = useRef(null);

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
                // blur the input
                curIdxRef.current.blur();
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
                    <input type="text" value={curIdx} className="input input-ghost input-xs text-center text-white input-primary w-10 font-mono" onChange={handleChange} ref={curIdxRef} />
                    <input type="text" value={index.max} className="input input-ghost input-xs text-center text-white input-primary w-10 font-mono disabled" readOnly />
                </div>

                <button className="btn btn-primary btn-sm text-2xl font-black" onClick={handles.next}>
                    <BsArrowBarRight />
                </button>
            </>
        )
    };

    const CurrentDirectory = () => {
        const [folders, setFolders] = useState([]);
        const [trashFolders, setTrashFolders] = useState([]);

        const fetchFolders = () => {
            axios.get('/api/folders')
                .then(res => {
                    const folders = res.data.folders;
                    folders.unshift('/');
                    setFolders(folders);
                })
                .catch(err => {
                    console.log(err);
                });
        }
        const fetchTrashFolders = () => {
            axios.get('/api/trash_folders')
                .then(res => {
                    const folders = res.data.folders;
                    folders.unshift('/');
                    setTrashFolders(folders);
                })
                .catch(err => {
                    console.log(err);
                });
        };

        useEffect(() => {
            fetchFolders();
            fetchTrashFolders();
        }, []);

        const handleFolderChange = (path, fromTrash) => {
            currentDir.set(path);
            currentDir.setFromTrash(fromTrash);

            const elem = document.activeElement;
            if (elem) {
                elem?.blur();
            }
        };

        return (
            <div className="flex-1 m-2 text-xl font-mono">Images:
                <div className="dropdown">
                    <label tabIndex={0} className="btn btn-ghost">
                        <div className='flex-1 overflow-clip'>
                            {currentDir.cur} {currentDir.fromTrash && <span className='text-red-500'> (trash)</span>}
                        </div>
                    </label>
                    <ul tabIndex={0} className="dropdown-content menu text-sm whitespace-nowrap p-2 bg-primary border-dashed w-auto">
                        {folders.map((folder) => {
                            return (
                                <li key={folder} className='hover:bg-accent' onClick={() => handleFolderChange(folder, false)}>
                                    <div className="flex flex-row py-0">
                                        <div className="flex-1 text-left">{folder}</div>
                                    </div>
                                </li>
                            )
                        })}
                        {trashFolders.map((folder) => {
                            return (
                                <li key={folder} className='hover:bg-accent' onClick={() => handleFolderChange(folder, true)}>
                                    <div className="flex flex-row py-0">
                                        <div className="flex-1 text-left">{folder}</div>
                                        <div className="flex-1 text-right text-red-500">trash</div>
                                    </div>
                                </li>
                            )
                        })}
                    </ul>
                </div>
            </div>
        )
    };

    return (
        <nav className="navbar bg-primary text-white py-0" ref={navRef}>
            <div className="flex flex-1">
                {CurrentDirectory()}
            </div>
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
