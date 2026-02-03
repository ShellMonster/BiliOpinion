import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ReportData } from '../../../types/report';

interface BrandNetworkProps {
  data: ReportData;
}

interface NetworkNode {
  id: string
  name: string
  value: number
  symbolSize: number
  category: number
  label: { show: boolean; fontWeight?: string }
  itemStyle: { color: string }
  tooltip?: { formatter: string }
}

interface NetworkLink {
  source: string
  target: string
}

/**
 * 品牌-型号关系网络图组件
 * 展示品牌与型号之间的从属关系，节点大小反映评论热度
 */
export const BrandNetwork: React.FC<BrandNetworkProps> = ({ data }) => {
  const option = useMemo(() => {
    // 提取数据
    const brands = data.brands || [];
    const modelRankings = data.model_rankings || [];
    const stats = data.stats;

    // 如果没有模型数据，只显示品牌或显示空提示
    if (modelRankings.length === 0) {
      return {
        title: {
          text: '品牌-型号关联网络',
          left: 'center'
        },
        series: []
      };
    }

    // 准备节点和边
    const nodes: NetworkNode[] = [];
    const links: NetworkLink[] = [];
    const addedNodes = new Set<string>();

    // 计算最大评论数用于归一化节点大小
    let maxComments = 0;
    if (stats && stats.comments_by_brand) {
      Object.values(stats.comments_by_brand).forEach(count => {
        maxComments = Math.max(maxComments, count);
      });
    }
    modelRankings.forEach(model => {
      maxComments = Math.max(maxComments, model.comment_count || 0);
    });
    
    // 避免除以零
    maxComments = maxComments || 100;

    // 辅助函数：计算节点大小
    const calculateSize = (count: number, isBrand: boolean) => {
      const baseSize = isBrand ? 30 : 10;
      const maxSize = isBrand ? 60 : 30;
      // 线性映射
      const size = baseSize + (count / maxComments) * (maxSize - baseSize);
      return Math.min(size, maxSize); // 限制最大尺寸
    };

    // 1. 添加品牌节点
    brands.forEach(brand => {
      if (!addedNodes.has(brand)) {
        const commentCount = stats?.comments_by_brand?.[brand] || 0;
        nodes.push({
          id: brand,
          name: brand,
          value: commentCount,
          symbolSize: calculateSize(commentCount, true),
          category: 0, // 品牌类别
          label: {
            show: true,
            fontWeight: 'bold'
          },
          itemStyle: {
            color: '#5470c6' // 品牌颜色
          }
        });
        addedNodes.add(brand);
      }
    });

    // 2. 添加型号节点和连接边
    modelRankings.forEach(model => {
      // 确保型号节点唯一 (有些型号名可能重复? 假设型号名在品牌下唯一，或者全局唯一)
      // 为了安全，可以使用 "Brand-Model" 作为ID，但显示 "Model"
      const modelId = `${model.brand}-${model.model}`;
      
      if (!addedNodes.has(modelId)) {
        nodes.push({
          id: modelId,
          name: model.model, // 显示名称
          value: model.comment_count,
          symbolSize: calculateSize(model.comment_count, false),
          category: 1, // 型号类别
          label: {
            show: true
          },
          itemStyle: {
            color: '#91cc75' // 型号颜色
          },
          // 可以在 tooltip 中显示更多信息
          tooltip: {
            formatter: `{b}<br/>所属品牌: ${model.brand}<br/>评论数: ${model.comment_count}<br/>评分: ${model.overall_score.toFixed(1)}`
          }
        });
        addedNodes.add(modelId);

        // 添加边：品牌 -> 型号
        if (addedNodes.has(model.brand)) {
          links.push({
            source: model.brand,
            target: modelId
          });
        }
      }
    });

    return {
      title: {
        text: '品牌-型号关联网络',
        subtext: '节点大小代表评论数量',
        left: 'center'
      },
      tooltip: {
        trigger: 'item'
      },
      legend: {
        data: ['品牌', '型号'],
        top: 'bottom'
      },
      series: [
        {
          type: 'graph',
          layout: 'force',
          data: nodes,
          links: links,
          categories: [
            { name: '品牌' },
            { name: '型号' }
          ],
          roam: true, // 允许缩放和平移
          label: {
            position: 'right',
            formatter: '{b}'
          },
          force: {
            repulsion: 300, // 节点之间的斥力因子
            edgeLength: [50, 150] // 边的两个节点之间的距离
          },
          lineStyle: {
            color: 'source',
            curveness: 0.1
          },
          emphasis: {
            focus: 'adjacency', // 高亮邻接节点
            lineStyle: {
              width: 4
            }
          }
        }
      ]
    };
  }, [data]);

  return (
    <div className="w-full bg-white p-4 rounded-lg shadow-sm">
      <ReactECharts 
        option={option} 
        style={{ height: '600px', width: '100%' }} 
        opts={{ renderer: 'canvas' }}
      />
    </div>
  );
};
