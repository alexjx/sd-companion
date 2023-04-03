
export default function Details(props) {
    const {
        file,
    } = props;

    return (
        <div className="flex flex-col text-white text-sm">
            <div className="">Name:</div>
            <div className="">{file?.path}</div>
            <div className="">Size:</div>
            <div className="">{file?.size}</div>
        </div>
    )
}