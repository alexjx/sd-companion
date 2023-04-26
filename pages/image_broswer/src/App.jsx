import { useState, useEffect, useRef } from 'react'
import axios from 'axios'
import { useHotkeys } from 'react-hotkeys-hook'

import Nav from './Components/Nav'
import Broswer from './Components/Broswer'
import PopupAlert from './Components/Alert'

import './App.css'

function App() {
    if (process.env.NODE_ENV === "development") {
        axios.defaults.baseURL = "http://localhost:9080";
    }

    const [currentDir, setCurrentDir] = useState("/")
    const [fromTrash, setFromTrash] = useState(false)
    const [files, setFiles] = useState([]);
    const [curIdx, setCurIdx] = useState(0);
    const [toJpeg, setToJpeg] = useState(true);

    const navRef = useRef(null);
    const containerRef = useRef(null);
    const [containerSize, setContainerSize] = useState({ height: 0, width: 0 });
    const imageRef = useRef(null);


    const handleNext = () => {
        if (curIdx + 1 < files.length) {
            setCurIdx((curIdx) => (curIdx + 1));
        } else {
            setCurIdx(0);
        }
    }

    const handlePrev = () => {
        if (curIdx > 0) {
            setCurIdx((curIdx) => (curIdx - 1));
        } else {
            setCurIdx(files.length - 1);
        }
    }

    const handleSetCurIdx = (idx) => {
        if (idx >= 0 && idx < files.length) {
            setCurIdx(idx);
        }
    }

    const handleDelete = () => {
        if (curIdx < files.length) {
            const f = files[curIdx];
            if (!fromTrash) {
                axios.delete(`/api/file?path=${f.path}`)
                    .then(res => {
                        fetchFiles();
                    })
            } else {
                axios.put(`/api/file?path=${f.path}`)
                    .then(res => {
                        fetchFiles();
                    })
            }
        }
    }

    function fetchFiles() {
        if (fromTrash) {
            axios.get(`/api/trash_files?dir=${currentDir}`)
                .then(res => {
                    setFiles(res.data.files);
                })
                .catch(err => {
                    // FIXME: alert user
                    console.log(err);
                });
        } else {
            axios.get(`/api/files?dir=${currentDir}`)
                .then(res => {
                    setFiles(res.data.files);
                })
                .catch(err => {
                    // FIXME: alert user
                    console.log(err);
                });
        }
    }

    // use effect to fetch issues
    useEffect(() => {
        fetchFiles();
    }, [currentDir, fromTrash]);
    useEffect(() => {
        if (containerRef.current === null) {
            return;
        }
        const updateSize = () => {
            setContainerSize({
                height: containerRef.current.clientHeight,
                width: containerRef.current.clientWidth,
            });
        };
        updateSize();
        window.addEventListener('resize', updateSize);
        return () => {
            window.removeEventListener('resize', updateSize);
        };
    }, [containerRef.current]);

    // hotkeys:
    //  left arrow: previous
    //  right arrow: next
    //  delete: delete
    //  r: refresh
    useHotkeys('left', handlePrev);
    useHotkeys('right', handleNext);
    useHotkeys('delete', handleDelete);
    useHotkeys('r', fetchFiles);
    useHotkeys('f', handleNext);
    useHotkeys('s', handlePrev)
    useHotkeys('shift+d', handleDelete);

    return (
        <div className='flex flex-col items-center bg-gray-800'>
            <Nav
                handles={
                    {
                        next: handleNext,
                        prev: handlePrev,
                        delete: handleDelete,
                        refresh: fetchFiles,
                        setCurIdx: handleSetCurIdx,
                    }
                }
                index={
                    {
                        cur: curIdx,
                        max: files.length,
                    }
                }
                currentDir={
                    {
                        cur: currentDir,
                        set: (dir) => {
                            setCurrentDir(dir);
                            setCurIdx(0);
                        },
                        fromTrash: fromTrash,
                        setFromTrash: (trash) => {
                            setFromTrash(trash);
                            setCurIdx(0);
                        },
                    }
                }
                containerSize={containerSize}
                jpegOpts={
                    {
                        toJpeg: toJpeg,
                        setToJpeg: setToJpeg,
                    }
                }
                navRef={navRef}
            />
            <Broswer
                files={files}
                root={axios.defaults.baseURL}
                curIdx={curIdx}
                innerRef={containerRef}
                toJpeg={toJpeg}
                imageRef={imageRef}
                fromTrash={fromTrash}
                containerSize={containerSize}
            />
        </div>
    )

}

export default App
