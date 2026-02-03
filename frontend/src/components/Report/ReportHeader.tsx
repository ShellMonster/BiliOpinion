
interface ReportHeaderProps {
  category: string
  reportId: string | number
}

const ReportHeader = ({
  category,
  reportId
}: ReportHeaderProps) => {
  return (
    <div className="flex justify-between items-end mb-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800">产品分析报告</h1>
        <p className="text-gray-500">{category} | 报告ID: {reportId}</p>
      </div>
    </div>
  )
}

export default ReportHeader
