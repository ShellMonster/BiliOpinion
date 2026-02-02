import React, { useState, useMemo } from 'react';
import type { Dimension, BrandRanking } from '../../types/report';

interface DecisionTreeProps {
  dimensions: Dimension[];
  rankings: BrandRanking[];
}

// 决策流程步骤
type Step = 'intro' | 'select-primary' | 'select-secondary' | 'result';

export const DecisionTree: React.FC<DecisionTreeProps> = ({ dimensions, rankings }) => {
  const [step, setStep] = useState<Step>('intro');
  const [primaryDimension, setPrimaryDimension] = useState<string>('');
  const [secondaryDimensions, setSecondaryDimensions] = useState<string[]>([]);

  // 计算推荐结果
  const recommendation = useMemo(() => {
    if (!primaryDimension) return null;

    // 评分权重: 主维度 60%, 次维度 40% (平分)
    const sortedBrands = [...rankings].sort((a, b) => {
      const getScore = (brand: BrandRanking, dim: string) => brand.scores[dim] || 0;

      let scoreA = getScore(a, primaryDimension) * 0.6;
      let scoreB = getScore(b, primaryDimension) * 0.6;

      if (secondaryDimensions.length > 0) {
        const secondaryWeight = 0.4 / secondaryDimensions.length;
        secondaryDimensions.forEach(dim => {
          scoreA += getScore(a, dim) * secondaryWeight;
          scoreB += getScore(b, dim) * secondaryWeight;
        });
      }

      return scoreB - scoreA;
    });

    return sortedBrands[0];
  }, [rankings, primaryDimension, secondaryDimensions]);

  const handleStart = () => setStep('select-primary');

  const handlePrimarySelect = (dimName: string) => {
    setPrimaryDimension(dimName);
    setStep('select-secondary');
  };

  const handleSecondaryToggle = (dimName: string) => {
    if (secondaryDimensions.includes(dimName)) {
      setSecondaryDimensions(prev => prev.filter(d => d !== dimName));
    } else {
      setSecondaryDimensions(prev => [...prev, dimName]);
    }
  };

  const handleFinish = () => setStep('result');

  const handleReset = () => {
    setStep('intro');
    setPrimaryDimension('');
    setSecondaryDimensions([]);
  };

  return (
    <div className="w-full max-w-2xl mx-auto p-6 bg-white rounded-xl shadow-lg border border-gray-100 min-h-[400px] flex flex-col relative overflow-hidden">
      {/* 背景装饰 */}
      <div className="absolute top-0 right-0 w-64 h-64 bg-blue-50 rounded-full -translate-y-1/2 translate-x-1/2 opacity-50 blur-3xl pointer-events-none" />
      <div className="absolute bottom-0 left-0 w-64 h-64 bg-purple-50 rounded-full translate-y-1/2 -translate-x-1/2 opacity-50 blur-3xl pointer-events-none" />

      <h2 className="text-2xl font-bold text-gray-800 mb-6 relative z-10 flex items-center gap-2">
        <span className="text-3xl">🤖</span> 
        <span>选购助手</span>
      </h2>

      <div className="flex-1 relative z-10">
          {step === 'intro' && (
            <div
              className="flex flex-col items-center justify-center h-full text-center space-y-6 py-8 animate-fade-in"
            >
              <div className="w-24 h-24 bg-blue-100 rounded-full flex items-center justify-center text-4xl mb-2 shadow-inner">
                🎯
              </div>
              <div>
                <h3 className="text-xl font-semibold text-gray-800 mb-2">纠结买哪个品牌？</h3>
                <p className="text-gray-500 max-w-md">
                  只需回答几个简单的问题，我们将根据全网评论数据，为您推荐最符合您需求的品牌。
                </p>
              </div>
              <button
                onClick={handleStart}
                className="px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-full font-medium transition-all transform hover:scale-105 shadow-lg shadow-blue-200"
              >
                开始挑选
              </button>
            </div>
          )}

          {step === 'select-primary' && (
            <div
              className="space-y-6 animate-fade-in"
            >
              <h3 className="text-lg font-medium text-gray-700">
                Q1. 在选购时，您<span className="text-blue-600 font-bold">最看重</span>哪一点？
              </h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {dimensions.map(dim => (
                  <button
                    key={dim.name}
                    onClick={() => handlePrimarySelect(dim.name)}
                    className="p-4 text-left border rounded-xl hover:border-blue-500 hover:bg-blue-50 transition-all group relative overflow-hidden"
                  >
                    <div className="font-semibold text-gray-800 group-hover:text-blue-700">{dim.name}</div>
                    <div className="text-sm text-gray-500 mt-1">{dim.description}</div>
                  </button>
                ))}
              </div>
            </div>
          )}

          {step === 'select-secondary' && (
            <div
              className="space-y-6 animate-fade-in"
            >
              <div>
                <h3 className="text-lg font-medium text-gray-700">
                  Q2. 还有其他<span className="text-blue-600 font-bold">关注的方面</span>吗？(可多选)
                </h3>
                <p className="text-sm text-gray-400 mt-1">已选择: {primaryDimension} (最重要)</p>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                {dimensions
                  .filter(d => d.name !== primaryDimension)
                  .map(dim => (
                    <button
                      key={dim.name}
                      onClick={() => handleSecondaryToggle(dim.name)}
                      className={`p-4 text-left border rounded-xl transition-all relative ${
                        secondaryDimensions.includes(dim.name)
                          ? 'border-blue-500 bg-blue-50 ring-1 ring-blue-500'
                          : 'border-gray-200 hover:bg-gray-50'
                      }`}
                    >
                      <div className="flex justify-between items-center">
                        <span className={`font-semibold ${secondaryDimensions.includes(dim.name) ? 'text-blue-700' : 'text-gray-800'}`}>
                          {dim.name}
                        </span>
                        {secondaryDimensions.includes(dim.name) && (
                          <span className="text-blue-500 text-lg">✓</span>
                        )}
                      </div>
                      <div className="text-sm text-gray-500 mt-1">{dim.description}</div>
                    </button>
                  ))}
              </div>

              <div className="flex justify-end pt-4">
                <button
                  onClick={handleFinish}
                  className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors shadow-md"
                >
                  查看推荐结果 →
                </button>
              </div>
            </div>
          )}

          {step === 'result' && recommendation && (
            <div
              className="space-y-6 text-center animate-fade-in"
            >
              <div className="inline-block px-4 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium mb-2">
                ✨ 为您匹配的最佳选择
              </div>
              
              <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-2xl p-8 border border-blue-100 shadow-sm relative overflow-hidden">
                <div className="absolute top-0 right-0 p-4 opacity-10">
                  <span className="text-9xl">🏆</span>
                </div>
                
                <h3 className="text-4xl font-bold text-gray-900 mb-2">{recommendation.brand}</h3>
                <div className="flex justify-center items-center gap-2 mb-6">
                  <span className="text-yellow-500 text-xl">★</span>
                  <span className="text-xl font-semibold text-gray-700">{recommendation.overall_score.toFixed(1)}</span>
                  <span className="text-gray-400 text-sm">(综合评分)</span>
                </div>

                <div className="space-y-3 bg-white/60 rounded-xl p-4 backdrop-blur-sm">
                  <div className="flex justify-between items-center p-2 border-b border-gray-100 last:border-0">
                    <span className="text-gray-600 font-medium">{primaryDimension} (首选)</span>
                    <span className="text-blue-600 font-bold text-lg">{recommendation.scores[primaryDimension]?.toFixed(1) || '-'}</span>
                  </div>
                  {secondaryDimensions.map(dim => (
                    <div key={dim} className="flex justify-between items-center p-2 border-b border-gray-100 last:border-0">
                      <span className="text-gray-600">{dim}</span>
                      <span className="text-gray-800 font-semibold">{recommendation.scores[dim]?.toFixed(1) || '-'}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex justify-center pt-4">
                <button
                  onClick={handleReset}
                  className="text-gray-500 hover:text-blue-600 flex items-center gap-2 transition-colors px-4 py-2 rounded-lg hover:bg-gray-50"
                >
                  <span>↺</span>
                  重新测试
                </button>
              </div>
            </div>
          )}
      </div>
    </div>
  );
};
