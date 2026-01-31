import { useParams } from 'react-router-dom'
import {
  Radar, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis,
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  Cell
} from 'recharts'

const Report = () => {
  const { id } = useParams()

  // Mock Data
  const reportData = {
    summary: {
      score: 8.5,
      recommendation: "强烈推荐",
      verdict: "该产品在充电速度和温控方面表现优异，用户口碑极佳。适合对性能有较高要求的用户，虽然价格略高，但物有所值。",
      pros: ["充电速度极快，实测满跑", "温控表现优秀，不烫手", "做工精致，手感好"],
      cons: ["价格相对较高", "体积略大，便携性一般"]
    },
    radarData: [
      { subject: '充电速度', A: 95, fullMark: 100 },
      { subject: '发热控制', A: 88, fullMark: 100 },
      { subject: '便携性', A: 70, fullMark: 100 },
      { subject: '做工质感', A: 90, fullMark: 100 },
      { subject: '性价比', A: 75, fullMark: 100 },
      { subject: '兼容性', A: 85, fullMark: 100 },
    ],
    sentimentData: [
      { name: '正面', value: 75, color: '#4ade80' },
      { name: '中立', value: 15, color: '#94a3b8' },
      { name: '负面', value: 10, color: '#f87171' },
    ],
    featureSatisfaction: [
      { name: '充电', score: 92 },
      { name: '散热', score: 85 },
      { name: '外观', score: 88 },
      { name: '价格', score: 72 },
      { name: '耐用', score: 80 },
    ]
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 space-y-6">
      <div className="flex justify-between items-end mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">产品分析报告</h1>
          <p className="text-gray-500">ID: {id} | 来源: Bilibili 评论分析</p>
        </div>
        <button className="px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors font-medium text-sm cursor-pointer">
          导出 PDF
        </button>
      </div>

      {/* Top Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Overall Score */}
        <div className="glass-card flex flex-col items-center justify-center p-8 bg-gradient-to-br from-blue-500 to-indigo-600 text-white border-none">
          <div className="text-6xl font-bold mb-2">{reportData.summary.score}</div>
          <div className="text-xl font-medium opacity-90">{reportData.summary.recommendation}</div>
          <div className="mt-4 px-3 py-1 bg-white/20 rounded-full text-sm backdrop-blur-sm">
            基于 1,243 条评论
          </div>
        </div>

        {/* Pros */}
        <div className="glass-card">
          <h3 className="text-lg font-semibold text-green-600 mb-4 flex items-center">
            <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
            </svg>
            核心优势
          </h3>
          <ul className="space-y-3">
            {reportData.summary.pros.map((item, i) => (
              <li key={i} className="flex items-start text-gray-700 text-sm">
                <span className="mr-2 text-green-500 mt-1">●</span>
                {item}
              </li>
            ))}
          </ul>
        </div>

        {/* Cons */}
        <div className="glass-card">
          <h3 className="text-lg font-semibold text-red-500 mb-4 flex items-center">
            <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.096c.5 0 .905-.405.905-.904 0-.715.211-1.413.608-2.008L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5" />
            </svg>
            待改进
          </h3>
          <ul className="space-y-3">
            {reportData.summary.cons.map((item, i) => (
              <li key={i} className="flex items-start text-gray-700 text-sm">
                <span className="mr-2 text-red-400 mt-1">●</span>
                {item}
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Main Analysis */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Radar Chart */}
        <div className="glass-card lg:col-span-1 min-h-[400px]">
          <h3 className="text-lg font-semibold text-gray-800 mb-4 text-center">六维能力图</h3>
          <div className="w-full h-[320px]">
            <ResponsiveContainer width="100%" height="100%">
              <RadarChart cx="50%" cy="50%" outerRadius="80%" data={reportData.radarData}>
                <PolarGrid stroke="#e5e7eb" />
                <PolarAngleAxis dataKey="subject" tick={{ fill: '#6b7280', fontSize: 12 }} />
                <PolarRadiusAxis angle={30} domain={[0, 100]} tick={false} axisLine={false} />
                <Radar
                  name="Product"
                  dataKey="A"
                  stroke="#3b82f6"
                  strokeWidth={3}
                  fill="#3b82f6"
                  fillOpacity={0.3}
                />
                <Tooltip />
              </RadarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Feature Satisfaction Bar Chart */}
        <div className="glass-card lg:col-span-2 min-h-[400px]">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">功能点满意度分析</h3>
          <div className="w-full h-[320px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={reportData.featureSatisfaction} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f0f0f0" />
                <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{ fill: '#6b7280' }} />
                <YAxis axisLine={false} tickLine={false} tick={{ fill: '#6b7280' }} />
                <Tooltip 
                  cursor={{ fill: '#f8fafc' }}
                  contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
                />
                <Bar dataKey="score" radius={[8, 8, 0, 0]} barSize={40}>
                  {reportData.featureSatisfaction.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.score > 80 ? '#3b82f6' : entry.score > 60 ? '#f59e0b' : '#ef4444'} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Buying Advice */}
      <div className="glass-card bg-gradient-to-r from-gray-50 to-gray-100 border border-gray-200">
        <h3 className="text-lg font-bold text-gray-800 mb-2">购买建议</h3>
        <p className="text-gray-600 leading-relaxed">
          {reportData.summary.verdict}
        </p>
      </div>
    </div>
  )
}

export default Report
