<!DOCTYPE html>
<html>
<head>
<style>
table, th, td {
  border: 1px solid black;
  border-collapse: collapse;
}
th, td {
  padding: 15px;
}
</style>
</head>
<body>

<h2>Table Users</h2>

<table style="width:100%">
  <tr>
                    <th>ID</th>
                    <th>Nama</th> 
                    <th>Umur</th>
                    <th>Alamat</th>
                    <th>Email</th>
                    <th>Role</th>
                </tr>
                <tr>
                    <td>{{.UserID}}</td>
                    <td>{{.Nama}}</td>
                    <td>{{.Umur}}</td>
                    <td>{{.Alamat}}</td>
                    <td>{{.Email}}</td>
                    <td>{{.Role}}</td>
                </tr>
</table>

</body>
</html>

