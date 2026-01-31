import { useParams } from 'react-router-dom'

export default function Progress() {
  const { id } = useParams()
  return (
    <div className="text-center py-20">
      <h1 className="text-3xl font-black text-slate-800 mb-4">
        分析进度: {id}
      </h1>
      <p className="text-slate-500">
        （进度页面将在Task 10实现）
      </p>
    </div>
  )
}
