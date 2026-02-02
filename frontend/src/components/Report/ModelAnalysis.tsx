import React, { useState, useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ModelRanking, Dimension } from '../../types/report';

interface ModelAnalysisProps {
  modelRankings: ModelRanking[];
  dimensions: Dimension[];
}

/**
 * 型号深度分析组件
 * 展示型号排名列表和雷达图对比
 */
export const ModelAnalysis: React.FC<ModelAnalysisProps> = ({ modelRankings, dimensions }) => {
  // 默认选中前3个型号进行对比
  const [selectedModels, setSelectedModels] = useState<string[]>(
    modelRankings.slice(0, 3).map(m => m.model)
  );

  // 切换型号选中状态
  const toggleModel = (model: string) => {
    if (selectedModels.includes(model)) {
      // 至少保留一个
      if (selectedModels.length > 1) {
        setSelectedModels(selectedModels.filter(m => m !== model));
      }
    } else {
      // 最多选择5个
      if (selectedModels.length < 5) {
        setSelectedModels([...selectedModels, model]);
      }
    }
  };

  // 预定义颜色盘
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#6366f1', '#ef4444', '#84cc16'];

  // 准备图表数据
  const chartOption = useMemo(() => {
    const selectedRankingData = modelRankings.filter(m => selectedModels.includes(m.model));
    
    return {
      title: {
        text: '型号性能对比',
        left: 'center',
        textStyle: {
          fontSize: 16,
          fontWeight: 'normal',
          color: '#374151'
        }
      },
      tooltip: {
        trigger: 'item'
      },
      legend: {
        bottom: 0,
        data: selectedRankingData.map(m => m.model),
        textStyle: {
          fontSize: 12
        }
      },
      radar: {
        indicator: dimensions.map(dim => ({
          name: dim.name,
          max: 100
        })),
        radius: '65%',
        center: ['50%', '50%'],
        splitNumber: 4,
        name: {
          textStyle: {
            color: '#6b7280',
            fontSize: 11
          }
        },
        splitLine: {
          lineStyle: {
            color: '#e5e7eb'
          }
        },
        splitArea: {
          show: true,
          areaStyle: {
            color: ['rgba(255, 255, 255, 0)', 'rgba(249, 250, 251, 0.5)']
          }
        }
      },
      series: [{
        type: 'radar',
        data: selectedRankingData.map((item, index) => {
          // 查找该型号在所有选中型号中的索引，用于分配颜色
          const colorIndex = index % colors.length;
          return {
            value: dimensions.map(dim => {
              const score = item.scores[dim.name] || 0;
              // 假设分数为0-10分制，如果是0-100则不需要乘以10
              // 根据BrandRadarChart逻辑，如果是小数需要放大
              return score <= 10 ? score * 10 : score;
            }),
            name: item.model,
            itemStyle: {
              color: colors[colorIndex]
            },
            lineStyle: {
              width: 2,
              color: colors[colorIndex]
            },
            areaStyle: {
              opacity: 0.1,
              color: colors[colorIndex]
            }
          };
        })
      }]
    };
  }, [modelRankings, dimensions, selectedModels, colors]);

  if (!modelRankings || modelRankings.length === 0) {
    return <div className="p-4 text-center text-gray-500">暂无型号分析数据</div>;
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
      <h2 className="text-lg font-bold text-gray-800 mb-6 flex items-center">
        <span className="w-1 h-6 bg-blue-500 rounded mr-2"></span>
        型号深度分析
      </h2>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* 左侧：排名列表 */}
        <div className="lg:col-span-7">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th scope="col" className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-10">
                    对比
                  </th>
                  <th scope="col" className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    排名
                  </th>
                  <th scope="col" className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    型号
                  </th>
                  <th scope="col" className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">
                    品牌
                  </th>
                  <th scope="col" className="px-3 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    综合评分
                  </th>
                  <th scope="col" className="px-3 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">
                    评论数
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {modelRankings.map((item) => (
                  <tr 
                    key={item.model} 
                    className={`hover:bg-gray-50 transition-colors ${selectedModels.includes(item.model) ? 'bg-blue-50/30' : ''}`}
                    onClick={() => toggleModel(item.model)}
                  >
                    <td className="px-3 py-4 whitespace-nowrap" onClick={(e) => e.stopPropagation()}>
                      <div className="flex items-center justify-center">
                        <input
                          type="checkbox"
                          className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded cursor-pointer"
                          checked={selectedModels.includes(item.model)}
                          onChange={() => toggleModel(item.model)}
                          disabled={!selectedModels.includes(item.model) && selectedModels.length >= 5}
                        />
                      </div>
                    </td>
                    <td className="px-3 py-4 whitespace-nowrap">
                      <div className={`flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold ${
                        item.rank <= 3 ? 'bg-yellow-100 text-yellow-700' : 'bg-gray-100 text-gray-600'
                      }`}>
                        {item.rank}
                      </div>
                    </td>
                    <td className="px-3 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">{item.model}</div>
                    </td>
                    <td className="px-3 py-4 whitespace-nowrap hidden sm:table-cell">
                      <div className="text-sm text-gray-500">{item.brand}</div>
                    </td>
                    <td className="px-3 py-4 whitespace-nowrap text-right">
                      <div className="text-sm font-bold text-gray-900">{item.overall_score.toFixed(1)}</div>
                    </td>
                    <td className="px-3 py-4 whitespace-nowrap text-right hidden sm:table-cell">
                      <div className="text-sm text-gray-500">{item.comment_count}</div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="mt-2 text-xs text-gray-500 text-right">
            * 点击行或复选框选择对比，最多选择5个型号
          </p>
        </div>

        {/* 右侧：雷达图 */}
        <div className="lg:col-span-5 flex flex-col">
          <div className="bg-gray-50 rounded-lg p-4 flex-grow flex items-center justify-center min-h-[350px]">
            <div className="w-full h-full">
              <ReactECharts 
                option={chartOption} 
                style={{ height: '350px', width: '100%' }} 
                opts={{ renderer: 'svg' }}
              />
            </div>
          </div>
          
          {/* 选中型号的详细维度得分 - 紧凑展示 */}
          <div className="mt-4 grid grid-cols-1 gap-2">
            {selectedModels.map((modelName, idx) => {
              const model = modelRankings.find(m => m.model === modelName);
              if (!model) return null;
              
              return (
                <div key={modelName} className="flex items-center text-xs border-b border-gray-100 pb-1 last:border-0">
                  <span 
                    className="w-3 h-3 rounded-full mr-2 flex-shrink-0" 
                    style={{ backgroundColor: colors[idx % colors.length] }}
                  ></span>
                  <span className="font-medium mr-2 truncate w-24" title={modelName}>{modelName}</span>
                  <div className="flex-grow flex justify-end gap-2 overflow-hidden">
                    {dimensions.slice(0, 3).map(dim => (
                      <span key={dim.name} className="text-gray-500 whitespace-nowrap">
                        {dim.name}: <span className="text-gray-900 font-medium">
                          {(model.scores[dim.name] || 0).toFixed(1)}
                        </span>
                      </span>
                    ))}
                    {dimensions.length > 3 && (
                      <span className="text-gray-400">...</span>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};
