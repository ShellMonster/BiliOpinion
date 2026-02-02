import React, { useState, useEffect } from 'react';
import type { BrandRanking, Dimension } from '../../types/report';

interface CompetitorCompareProps {
  rankings: BrandRanking[];
  dimensions: Dimension[];
}

export const CompetitorCompare: React.FC<CompetitorCompareProps> = ({ rankings, dimensions }) => {
  // 初始化选择的品牌，默认选择前两名
  const [brand1Name, setBrand1Name] = useState<string>('');
  const [brand2Name, setBrand2Name] = useState<string>('');

  useEffect(() => {
    if (rankings.length > 0) {
      if (!brand1Name) setBrand1Name(rankings[0].brand);
      if (!brand2Name && rankings.length > 1) setBrand2Name(rankings[1].brand);
      else if (!brand2Name && rankings.length > 0) setBrand2Name(rankings[0].brand);
    }
  }, [rankings]);

  // 获取选中品牌的数据
  const brand1 = rankings.find(r => r.brand === brand1Name);
  const brand2 = rankings.find(r => r.brand === brand2Name);

  if (!brand1 || !brand2) {
    return null;
  }

  // 判断优势方
  const getWinner = (score1: number, score2: number) => {
    if (score1 > score2) return 'brand1';
    if (score2 > score1) return 'brand2';
    return 'tie';
  };

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="p-6 border-b border-gray-100">
        <h2 className="text-xl font-bold text-gray-800 flex items-center gap-2">
          <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
          </svg>
          竞品深度对比 (1v1)
        </h2>
        <p className="text-gray-500 text-sm mt-1">选择两个品牌进行详细的维度得分对比分析</p>
      </div>

      <div className="p-6">
        {/* 品牌选择器区域 */}
        <div className="grid grid-cols-[1fr_auto_1fr] gap-4 items-center mb-8">
          {/* 品牌 1 选择 */}
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium text-gray-600">选择品牌 A</label>
            <select
              value={brand1Name}
              onChange={(e) => setBrand1Name(e.target.value)}
              className="w-full p-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all font-medium text-gray-800"
            >
              {rankings.map((r) => (
                <option key={r.brand} value={r.brand}>
                  {r.brand} (第{r.rank}名)
                </option>
              ))}
            </select>
          </div>

          {/* VS 徽章 */}
          <div className="flex flex-col items-center justify-center pt-6">
            <div className="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center font-black text-gray-400 text-sm italic border-2 border-white shadow-sm">
              VS
            </div>
          </div>

          {/* 品牌 2 选择 */}
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium text-gray-600 text-right">选择品牌 B</label>
            <select
              value={brand2Name}
              onChange={(e) => setBrand2Name(e.target.value)}
              className="w-full p-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all font-medium text-gray-800 text-right"
              dir="rtl"
            >
              {rankings.map((r) => (
                <option key={r.brand} value={r.brand}>
                  (第{r.rank}名) {r.brand}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* 核心指标对比 */}
        <div className="grid grid-cols-3 gap-4 mb-8 bg-gray-50 rounded-xl p-4">
          <div className="text-center">
            <div className={`text-2xl font-bold ${brand1.overall_score > brand2.overall_score ? 'text-blue-600' : 'text-gray-700'}`}>
              {brand1.overall_score.toFixed(1)}
            </div>
            <div className="text-xs text-gray-500 mt-1">综合评分</div>
            {brand1.overall_score > brand2.overall_score && (
              <div className="inline-flex items-center gap-1 mt-1 text-xs font-medium text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full">
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
                </svg> 胜出
              </div>
            )}
          </div>
          
          <div className="flex items-center justify-center text-gray-400 text-sm font-medium">
            综合表现
          </div>

          <div className="text-center">
            <div className={`text-2xl font-bold ${brand2.overall_score > brand1.overall_score ? 'text-blue-600' : 'text-gray-700'}`}>
              {brand2.overall_score.toFixed(1)}
            </div>
            <div className="text-xs text-gray-500 mt-1">综合评分</div>
            {brand2.overall_score > brand1.overall_score && (
              <div className="inline-flex items-center gap-1 mt-1 text-xs font-medium text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full">
                 <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
                </svg> 胜出
              </div>
            )}
          </div>
        </div>

        {/* 详细维度对比 */}
        <div className="space-y-6">
          <h3 className="text-sm font-bold text-gray-900 uppercase tracking-wider mb-4 border-l-4 border-blue-500 pl-3">维度详细对比</h3>
          
          {dimensions.map((dim) => {
            const score1 = brand1.scores[dim.name] || 0;
            const score2 = brand2.scores[dim.name] || 0;
            const winner = getWinner(score1, score2);
            
            return (
              <div key={dim.name} className="relative group">
                <div className="flex justify-between items-end mb-2 text-sm">
                  <span className={`font-medium ${winner === 'brand1' ? 'text-blue-700' : 'text-gray-600'}`}>
                    {score1.toFixed(1)}
                  </span>
                  <div className="flex flex-col items-center">
                    <span className="font-bold text-gray-800">{dim.name}</span>
                    <span className="text-[10px] text-gray-400 max-w-[150px] truncate text-center hidden group-hover:block absolute -top-4 bg-gray-800 text-white px-2 py-0.5 rounded transition-all">
                      {dim.description}
                    </span>
                  </div>
                  <span className={`font-medium ${winner === 'brand2' ? 'text-blue-700' : 'text-gray-600'}`}>
                    {score2.toFixed(1)}
                  </span>
                </div>
                
                {/* 进度条对比 */}
                <div className="flex items-center gap-2 h-2.5">
                  {/* Brand 1 Bar (Right aligned) */}
                  <div className="flex-1 flex justify-end bg-gray-100 rounded-l-full overflow-hidden h-full">
                    <div 
                      className={`h-full rounded-l-full transition-all duration-500 ${winner === 'brand1' ? 'bg-blue-500' : 'bg-gray-300'}`}
                      style={{ width: `${(score1 / 10) * 100}%` }}
                    ></div>
                  </div>
                  
                  {/* Center Divider */}
                  <div className="w-px h-4 bg-gray-200"></div>

                  {/* Brand 2 Bar (Left aligned) */}
                  <div className="flex-1 flex justify-start bg-gray-100 rounded-r-full overflow-hidden h-full">
                    <div 
                      className={`h-full rounded-r-full transition-all duration-500 ${winner === 'brand2' ? 'bg-blue-500' : 'bg-gray-300'}`}
                      style={{ width: `${(score2 / 10) * 100}%` }}
                    ></div>
                  </div>
                </div>

                {/* 差异说明 */}
                {(score1 !== score2) && (
                  <div className="text-center mt-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    <span className="text-[10px] text-gray-400 bg-gray-50 px-2 py-0.5 rounded-full border border-gray-100">
                      {score1 > score2 ? `${brand1.brand}` : `${brand2.brand}`} +{Math.abs(score1 - score2).toFixed(1)}
                    </span>
                  </div>
                )}
              </div>
            );
          })}
        </div>
        
        {/* 底部提示 */}
        <div className="mt-8 pt-4 border-t border-gray-100 flex items-start gap-2 text-xs text-gray-400">
          <svg className="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p>对比数据基于AI对评论内容的语义分析，分数范围为1-10分。高分代表在该维度上用户评价更正面。</p>
        </div>
      </div>
    </div>
  );
};