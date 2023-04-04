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
                setMetadata(res.data.metadata);
            })
            .catch(err => {
                console.log(err);
            });
    };

    useEffect(() => {
        if (file) {
            fetchMetadata(file.path);
        }
    }, [file]);

    const formatOptions = (sd_info) => {
        let optKeys = Object.keys(sd_info.options);
        optKeys.sort();
        return (
            <>
                {optKeys.map((key) => {
                    return (
                        <div className="grid grid-cols-5" key={key}>
                            <div className="col-span-2 text-gray-200">{key}</div>
                            <div className="col-span-3">{sd_info.options[key]}</div>
                        </div>
                    )
                })}
            </>
        )
    }

    const formatLoraInfo = (sd_info) => {
        if (sd_info.loras) {
            // sort by name, item in {name: xxx, weight: xxx}
            sd_info.loras.sort((a, b) => {
                return a.name.localeCompare(b.name);
            });

            let totalWeight = 0;

            const weightColor = (weight) => {
                if (weight < 1.0) {
                    return "text-white";
                } else if (weight < 1.5) {
                    return "text-yellow-300";
                } else {
                    return "text-red-300";
                }
            }

            return (
                <div className="flex flex-col">
                    <div className="text-gray-200 text-base my-2">LORA:</div>
                    {
                        sd_info.loras.map((lora) => {
                            totalWeight += lora.weight;
                            return (
                                <div className="grid grid-cols-4" key={lora.name}>
                                    <div className="text-gray-200 col-span-3">{lora.name}</div>
                                    <div className="col-span-1">{lora.weight}</div>
                                </div>
                            )
                        })
                    }
                    <div className="grid grid-cols-4 text-gray-400">
                        <div className="col-span-3">Total</div>
                        <div className={"col-span-1 "+weightColor(totalWeight)}>{totalWeight.toFixed(1)}</div>
                    </div>
                </div>
            )
        }
    }

    const formatSDInfo = (meta) => {
        if (meta.sd_params) {
            return (
                <>
                    {formatOptions(metadata.sd_params)}
                    <div className="grid grid-cols-1">
                        <div className="text-gray-200 text-base my-2">Prompt:</div>
                        <div>{metadata.sd_params.prompt}</div>
                        {
                            metadata.sd_params.negative_prompt && (
                                <>
                                    <div className="text-gray-200 text-base my-2">Negative Prompt:</div>
                                    <div>{metadata.sd_params.negative_prompt}</div>
                                </>
                            )
                        }
                        {
                            // metadata.sd_params.template && (
                            //     <>
                            //         <div className="text-gray-200 text-base my-2">Template:</div>
                            //         <div className="whitespace-pre-wrap">{metadata.sd_params.template}</div>
                            //     </>
                            // )
                        }
                        {
                           formatLoraInfo(metadata.sd_params)
                        }
                    </div>
                </>
            )
        }
    }

    const handleCopyPropmpt = () => {
        // reconsturct the prompt
        let prompt = metadata.sd_params.prompt;
        if (metadata.sd_params.negative_prompt) {
            prompt += "\nNegative prompt: " + metadata.sd_params.negative_prompt;
        }
        if (metadata.sd_params.option_str) {
            prompt += "\n" + metadata.sd_params.option_str;
        }
        navigator.clipboard.writeText(prompt);
    }

    return (
        <div className="flex flex-col text-white text-xs">
            {window.isSecureContext && metadata.sd_params && (
                <button className="bg-gray-800 text-gray-300 rounded-md px-2 py-1 my-2" onClick={handleCopyPropmpt}>Copy Prompt</button>
            )}
            <div className="grid grid-cols-5">
                <div className="col-span-2 text-gray-300">Name</div>
                <div className="col-span-3">
                    <a href={file && `/files/${file?.path}`}>{metadata.name}</a>
                </div>
                <div className="col-span-2 text-gray-300">Size</div>
                <div classname="col-span-3">
                    {metadata.size}
                </div>
            </div>
            {formatSDInfo(metadata)}
        </div>
    )
}
