import 'package:admin/shoppinglist_page.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:admin/models/food_types.dart';
import 'package:admin/foods_page.dart';
import 'package:flutter/foundation.dart';
import 'package:file_picker/file_picker.dart';

class FoodTypePage extends StatefulWidget {
  const FoodTypePage({super.key});

  @override
  State<FoodTypePage> createState() => _FoodTypePageState();
}

class _FoodTypePageState extends State<FoodTypePage> {
  PlatformFile? _pickedFile;


  Future<List<FoodType>> fetchFoodTypes() async {
    final String url = kIsWeb 
      ? '/itemtypes' 
      : 'http://192.168.10.60/itemtypes';

    final response = await http.get(Uri.parse(url));

    if (response.statusCode == 200) {
      List jsonResponse = json.decode(response.body);
      return jsonResponse.map((data) => FoodType.fromJson(data)).toList();
    } else {
      throw Exception('Could not load food types');
    }
  }

  Future<void> _addNewFoodType(String name, PlatformFile pickedFile) async {
  final String url = kIsWeb 
    ? '/itemtypes' 
    : 'http://192.168.10.60/itemtype';

  var request = http.MultipartRequest('POST', Uri.parse(url));
  request.headers.addAll({
    "Accept": "*/*",
    "Content-Type": "multipart/form-data",
  });
  request.fields['name'] = name;

  if (pickedFile.bytes != null) {
    request.files.add(http.MultipartFile.fromBytes(
      'image', 
      pickedFile.bytes!,
      filename: pickedFile.name,
    ));
  } else {
    print("Hiba: A fájl bájtok üresek!");
    return;
  }

  var streamedResponse = await request.send();
  var response = await http.Response.fromStream(streamedResponse);

  if (response.statusCode == 201 || response.statusCode == 200) {
    setState(() {});
    
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Sikeresen hozzáadva!')),
    );
  } else {
    throw Exception('Sikertelen mentés: ${response.body}');
  }
}

  void _showAddDialog() {
  TextEditingController nameController = TextEditingController();
  PlatformFile? localPickedFile;

  showDialog(
    context: context,
    builder: (context) => StatefulBuilder(
      builder: (context, setDialogState) => AlertDialog(
        title: const Text('Új étel típus'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController, 
              decoration: const InputDecoration(labelText: 'Étel neve')
            ),
            const SizedBox(height: 20),
            ElevatedButton.icon(
              onPressed: () async {
                FilePickerResult? result = await FilePicker.platform.pickFiles(
                  type: FileType.image,
                  allowMultiple: false,
                  withData: true,
                );

                if (result != null) {
                  setDialogState(() {
                    localPickedFile = result.files.first;
                  });
                }
              },
              icon: const Icon(Icons.image),
              label: Text(localPickedFile == null 
                  ? 'Kép kiválasztása' 
                  : 'Kép: ${localPickedFile!.name}'),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Mégse')),
          ElevatedButton(
            onPressed: (localPickedFile == null) ? null : () {
              _addNewFoodType(nameController.text, localPickedFile!);
              Navigator.pop(context);
            }, 
            child: const Text('Mentés')
          ),
        ],
      ),
    ),
  );
}

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Food Types'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      drawer: Drawer(
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            const DrawerHeader(
              decoration: BoxDecoration(color: Colors.deepPurple),
              child: Text('Menu', style: TextStyle(color: Colors.white, fontSize: 24)),
            ),
            ListTile(
              leading: const Icon(Icons.add_task),
              title: const Text('Shoppinglist'),
              onTap: () {
                Navigator.pop(context);
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ShoppingListPage()),
                );
              },
            ),
            ListTile(
              leading: const Icon(Icons.add_shopping_cart),
              title: const Text('Inventory'),
              onTap: () {
                Navigator.pop(context);
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => FoodPage()),
                );
              },
            ),
            ListTile(
              leading: const Icon(Icons.fastfood),
              title: const Text('Food types'),
              onTap: () {
                Navigator.pop(context);
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => FoodTypePage()),
                );
              },
            ),
          ],
        ),
      ),
      body: FutureBuilder<List<FoodType>>(
        future: fetchFoodTypes(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Hiba: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('Nincs adat.'));
          }

          return ListView.builder(
            itemCount: snapshot.data!.length,
            itemBuilder: (context, index) {
              final item = snapshot.data![index];
              return ListTile(
                leading: const CircleAvatar(child: Icon(Icons.restaurant)),
                title: Text(item.name),
                subtitle: Text('ID: ${item.id}'),
                onTap: () {
                   print('Selected: ${item.name}');
                },
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showAddDialog,
        backgroundColor: Colors.deepPurple,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }
}