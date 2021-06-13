import React from "react";
import "./style.css";

const UploadMultipart = () => {
  const nameInput = React.useRef<HTMLInputElement>(null);
  const fileInput = React.useRef<HTMLInputElement>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (
      nameInput?.current?.value &&
      fileInput?.current?.files &&
      fileInput?.current?.files[0]
    ) {
      const formData = new FormData();
      formData.append("id", nameInput.current.value);
      formData.append("file", fileInput.current.files[0]);

      fetch(process.env.REACT_APP_IMAGES_ENDPOINT + "/", {
        method: "POST",
        headers: {
          // Remove that header and allow `fetch` to generate the full content-type
          // @see https://github.com/github/fetch/issues/505#issuecomment-293064470
          // 'Content-Type': `multipart/form-data; boundary=${formData._boundary}`
        },
        body: formData,
      });
    }
  };

  return (
    <div className="upload-multipart">
      <h3>Upload</h3>
      <form onSubmit={handleSubmit}>
        <label>
          ID:
          <input ref={nameInput} type="number" name="id" />
        </label>
        <br />
        <label>
          File:
          <input ref={fileInput} type="file" name="file" />
        </label>
        <input type="submit" value="Submit" />
      </form>
    </div>
  );
};

export default UploadMultipart;
