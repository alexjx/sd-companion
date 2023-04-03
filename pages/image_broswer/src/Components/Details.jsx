import { useEffect, useState } from "react";
import axios from "axios";

export default function Details(props) {
    const {
        file,
    } = props;

    const [metadata, setMetadata] = useState([]);

    const fetchMetadata = (path) => {
        axios.get(`/api/metadata?path=${path}`)
            .then(res => {
                let metadata = [];
                let prompt, negative, template, name, size;

                // sort keys
                let keys = Object.keys(res.data.metadata);
                keys.sort();

                // map keys
                keys.map((key, i) => {
                    if (key === "Prompt") {
                        prompt = res.data.metadata[key];
                    } else if (key === "Negative prompt") {
                        negative = res.data.metadata[key];
                    } else if (key === "Template") {
                        template = res.data.metadata[key];
                    } else if (key === "Name") {
                        name = res.data.metadata[key];
                    } else if (key === "Size") {
                        size = res.data.metadata[key];
                    } else {
                        metadata.push({
                            key: key,
                            value: res.data.metadata[key],
                            type: 1, // inline
                        });
                    }
                });

                if (name) {
                    metadata.splice(0, 0, {
                        key: "Name",
                        value: name,
                        type: 1, // inline
                    });
                }
                if (size) {
                    metadata.splice(1, 0, {
                        key: "Size",
                        value: size,
                        type: 1, // inline
                    });
                }
                if (prompt) {
                    metadata.push({
                        key: "Prompt",
                        value: prompt,
                        type: 2, // block
                    });
                }
                if (negative) {
                    metadata.push({
                        key: "Negative prompt",
                        value: negative,
                        type: 2, // block
                    });
                }
                if (template) {
                    // tempalte has newlines, we need to preserve
                    metadata.push({
                        key: "Template",
                        value: template,
                        type: 3, // block
                    });
                }

                setMetadata(metadata);
            })
            .catch(err => {
                console.log(err);
            });
    };

    useEffect(() => {
        fetchMetadata(file.path);
    }, [file]);

    const formatMetadata = (metadata) => {
        // map keys for div pairs
        return metadata.map((item, i) => {
            let className;
            switch (item.type) {
                case 1:
                    className = "grid grid-cols-2";
                    break;
                case 2:
                    className = "flex flex-col";
                    break;
                case 3:
                    className = "flex flex-col whitespace-pre-wrap";
                    break;
                default:
                    className = "grid grid-cols-2";
                    break;
            }
            return (
                <div key={i} className={className}>
                    <div className="mr-5">{item.key}:</div>
                    <div>{item.value}</div>
                </div>
            )
        })
    };

    return (
        <div className="flex flex-col text-white text-xs">
            {formatMetadata(metadata)}
        </div>
    )
}
