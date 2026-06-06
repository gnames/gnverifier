$(function () {
  var $table = $('#data-sources');

  $('#ds-filter').on('input', function () {
    var q = $(this).val().trim().toLowerCase();
    var visibleIdx = 0;
    $table.find('tr.even, tr.odd').each(function () {
      var $tds = $(this).find('td');
      var title = $tds.eq(1).text().toLowerCase();
      var shortTitle = $tds.eq(2).text().toLowerCase();
      var match = !q || title.includes(q) || shortTitle.includes(q);
      $(this).toggle(match);
      if (match) {
        $(this).removeClass('even odd').addClass(visibleIdx % 2 === 0 ? 'even' : 'odd');
        visibleIdx++;
      }
    });
  });

  var sortCol = -1;
  var sortAsc = true;

  var $headers = $table.find('th').css('cursor', 'pointer');
  $headers.first().append(' <span class="sort-indicator">▲</span>');
  sortCol = 0;

  $headers.each(function (colIdx) {
    $(this).on('click', function () {
      if (sortCol === colIdx) {
        sortAsc = !sortAsc;
      } else {
        sortCol = colIdx;
        sortAsc = true;
      }

      $table.find('th').each(function (i) {
        $(this).find('.sort-indicator').remove();
        if (i === sortCol) {
          $(this).append(' <span class="sort-indicator">' + (sortAsc ? '▲' : '▼') + '</span>');
        }
      });

      var $tbody = $table.find('tr.even, tr.odd');
      var rows = $tbody.toArray();

      rows.sort(function (a, b) {
        var aText = $(a).find('td').eq(colIdx).text().trim();
        var bText = $(b).find('td').eq(colIdx).text().trim();
        var aNum = Number(aText);
        var bNum = Number(bText);
        var cmp;
        if (!isNaN(aNum) && !isNaN(bNum)) {
          cmp = aNum - bNum;
        } else {
          cmp = aText.localeCompare(bText);
        }
        return sortAsc ? cmp : -cmp;
      });

      $.each(rows, function (i, row) {
        $(row).removeClass('even odd').addClass(i % 2 === 0 ? 'even' : 'odd');
        $table.append(row);
      });
    });
  });
});
