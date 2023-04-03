import { useState, useEffect, useRef } from 'react'
import axios from 'axios'
import Nav from './Components/Nav'
import Broswer from './Components/Broswer'
import { useHotkeys } from 'react-hotkeys-hook'

function App() {
    if (process.env.NODE_ENV === "development") {
        axios.defaults.baseURL = "http://localhost:9080";
    }

    const [rootPath, setRootPath] = useState('');
    const [files, setFiles] = useState([]);
    const [curIdx, setCurIdx] = useState(0);

    const containerRef = useRef(null);
    const [containerSize, setContainerSize] = useState({ height: 0, width: 0 });

    const handleNext = () => {
        console.log('next', curIdx);
        if (curIdx + 1 < files.length) {
            setCurIdx((curIdx) => (curIdx + 1));
        } else {
            setCurIdx(0);
        }
    }

    const handlePrev = () => {
        console.log('prev', curIdx);
        if (curIdx > 0) {
            setCurIdx((curIdx) => (curIdx - 1));
        } else {
            setCurIdx(files.length - 1);
        }
    }

    const handleDelete = () => {
        if (curIdx < files.length) {
            const f = files[curIdx];
            axios.delete(`/api/file?path=${f.path}`)
            .then(res => {
                fetchFiles();
            })
        }
    }

    function fetchRootPath() {
        axios.get(`/api/root`)
            .then(res => {
                setRootPath(res.data.root);
            })
            .catch(err => {
                setRootPath('');
            });
    }

    function fetchFiles() {
        axios.get(`/api/files`)
            .then(res => {
                setFiles(res.data.files);
            })
            .catch(err => {
                // FIXME: alert user
                console.log(err);
            });
    }

    // use effect to fetch issues
    useEffect(() => {
        fetchFiles();
    }, []);
    useEffect(() => {
        fetchRootPath();
    }, []);
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

    return (
        <div className='flex flex-col h-screen'>
            <Nav
                handles={
                    {
                        next: handleNext,
                        prev: handlePrev,
                        delete: handleDelete,
                        refresh: fetchFiles,
                    }
                }
                index={
                    {
                        cur: curIdx,
                        max: files.length,
                    }
                }
                rootPath={rootPath}
                containerSize={containerSize}
            />
            <Broswer
                files={files}
                root={axios.defaults.baseURL}
                curIdx={curIdx}
                innerRef={containerRef}
                containerSize={containerSize}
            />
        </div>
    )

}

export default App
