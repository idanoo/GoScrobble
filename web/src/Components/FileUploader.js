import React from 'react'

const FileUploader = ({onFileSelect}) => {
    const handleFileInput = (e) => {
        // handle validations here in future
        onFileSelect(e.target.files[0]);
      };

    return (
        <span>
            <input type="file" onChange={handleFileInput} />
        </span>
    )
}

export default FileUploader