import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ReportData } from '../../../types/report';

interface BrandHeatmapProps {
  data: ReportData;
}

interface HeatmapTooltipParams {
  value: [number, number, number]
}

/**
 * 品牌维度得分热力图组件
 * 展示品牌在各个维度上的得分表现，使用颜色区分得分等级
 */
export const BrandHeatmap: React.FC<BrandHeatmapProps> = ({ data }) => {
  const option = useMemo(() => {
    // 准备 X 轴（维度）和 Y 轴（品牌）数据
    const dimensions = data.dimensions.map(d => d.name);
    const brands = data.brands;

    // 构建热力图数据: [dimensionIndex, brandIndex, score]
    // x: 维度, y: 品牌
    const seriesData = [];
    for (let i = 0; i < brands.length; i++) {
      const brand = brands[i];
      for (let j = 0; j < dimensions.length; j++) {
        const dimension = dimensions[j];
        // 获取得分，默认为 0
        const score = data.scores[brand]?.[dimension] || 0;
        // 只有当有分数时才显示
        if (data.scores[brand]?.[dimension] !== undefined) {
          seriesData.push([j, i, score]);
        }
      }
    }

    return {
      tooltip: {
        position: 'top',
        formatter: (params: HeatmapTooltipParams) => {
          const dimIndex = params.value[0];
          const brandIndex = params.value[1];
          const score = params.value[2];
          return `
            <div class="font-bold">${brands[brandIndex]}</div>
            <div>${dimensions[dimIndex]}: <span class="font-bold">${score}</span></div>
          `;
        }
      },
      grid: {
        height: '65%',
        top: '10%',
        right: '5%',
        left: '22%' // 增加空间给品牌名称，避免截断
      },
      xAxis: {
        type: 'category',
        data: dimensions,
        splitArea: {
          show: true
        },
        axisLabel: {
          interval: 0,
          rotate: 30 // 防止维度名称重叠
        }
      },
      yAxis: {
        type: 'category',
        data: brands,
        splitArea: {
          show: true
        }
      },
      visualMap: {
        min: 0,
        max: 10,
        calculable: true,
        orient: 'horizontal',
        left: 'center',
        bottom: '0%',
        type: 'piecewise', // 分段型
        pieces: [
          { gte: 8, label: '优秀 (≥8)', color: '#10b981' }, // 绿色
          { gte: 6, lt: 8, label: '良好 (6-8)', color: '#f59e0b' }, // 橙色
          { lt: 6, label: '需改进 (<6)', color: '#ef4444' } // 红色
        ],
        textStyle: {
          color: '#374151'
        }
      },
      series: [
        {
          name: '品牌得分',
          type: 'heatmap',
          data: seriesData,
          label: {
            show: true,
            color: '#374151'
          },
          itemStyle: {
            borderColor: '#fff',
            borderWidth: 1
          },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    };
  }, [data]);

  return (
    <div className="w-full bg-white rounded-lg p-4 shadow-sm border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">品牌维度热力图</h3>
      <ReactECharts option={option} style={{ height: '400px', width: '100%' }} />
    </div>
  );
};
