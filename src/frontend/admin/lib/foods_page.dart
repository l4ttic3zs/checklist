import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'models/foods.dart';

class FoodTypePage extends StatefulWidget {
  const FoodTypePage({super.key});

  @override
  State<FoodTypePage> createState() => _FoodTypePageState();
}

class _FoodTypePageState extends State<FoodTypePage> {
  Future<List<FoodType>> fetchFoodTypes() async {
    final response = await http.get(Uri.parse('http://http://192.168.10.60/itemtypes'));

    if (response.statusCode == 200) {
      List jsonResponse = json.decode(response.body);
      return jsonResponse.map((data) => FoodType.fromJson(data)).toList();
    } else {
      throw Exception('Could not load food types');
    }
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
              },
            ),
            ListTile(
              leading: const Icon(Icons.add_shopping_cart),
              title: const Text('Inventory'),
              onTap: () {
                Navigator.pop(context);
              },
            ),
            ListTile(
              leading: const Icon(Icons.fastfood),
              title: const Text('Food types'),
              onTap: () {
                Navigator.pop(context);
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
    );
  }
}