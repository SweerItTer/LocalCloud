import { useState } from 'react'

function App() {
  const [uploading, setUploading] = useState(false)
  const [preview, setPreview] = useState<string>()

  const handleUpload = async () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'image/*'
    
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0]
      if (!file) return

      setUploading(true)
      
      try {
        const formData = new FormData()
        formData.append('file', file)
        
        const res = await fetch('/api/upload', {
          method: 'POST',
          body: formData
        })
        
        const data = await res.json()
        if (data.url) {
          // 显示预览图
          setPreview(`http://localhost:8080/api/image/${data.url.split('/images/')[1]}`)
        }
      } catch (err) {
        alert('上传失败')
      } finally {
        setUploading(false)
      }
    }

    input.click()
  }

  return (
    <div style={{ padding: 20 }}>
      <button 
        onClick={handleUpload}
        disabled={uploading}
        style={{ fontSize: 18, padding: 10 }}
      >
        {uploading ? '上传中...' : '选择图片上传'}
      </button>
      
      {preview && (
        <div style={{ marginTop: 20 }}>
          <img 
            src={preview} 
            alt="预览" 
            style={{ maxWidth: '100%', maxHeight: 400 }}
          />
        </div>
      )}
    </div>
  )
}

export default App