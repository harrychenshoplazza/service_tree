document.addEventListener('DOMContentLoaded', () => {
    // 树形结构交互
    document.querySelectorAll('.tree-node li').forEach(node => {
        if (node.querySelector('ul')) {
            node.classList.add('has-child');
            node.style.cursor = 'pointer';
        }

        node.addEventListener('click', (e) => {
            e.stopPropagation();
            const ul = node.querySelector('ul');
            if (ul) {
                const isHidden = ul.style.display === 'none';
                ul.style.display = isHidden ? 'block' : 'none';
                node.classList.toggle('expanded', isHidden);
            }
        });
    });

    // 详情按钮交互
    document.querySelectorAll('.detail-button').forEach(button => {
        button.addEventListener('click', (e) => {
            e.stopPropagation();
            const card = button.closest('.server-card');
            const ip = card.querySelector('.server-field span').nextSibling.textContent.trim();
            alert(`服务器详情：\nIP: ${ip}`); // 示例：显示 IP 详情
        });
    });
});