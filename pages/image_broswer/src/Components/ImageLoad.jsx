import { useState } from 'react';
import { MdOutlineDownloading } from 'react-icons/md';

const ImageLoader = ({ src, imgRef }) => {
    const [loaded, setLoaded] = useState(false);

    // Define a function to handle the image load event
    const handleImageLoad = () => {
        setLoaded(true);
    };

    return (
        <div className="h-full-important">
            <img
                key={src}
                src={src}
                className={`${loaded ? '' : 'hidden'} h-full-important object-contain`}
                ref={imgRef}
                onLoad={handleImageLoad}
            />
        </div>
    );
};

export default ImageLoader;