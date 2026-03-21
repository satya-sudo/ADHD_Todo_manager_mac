#import <Foundation/Foundation.h>
#import <UserNotifications/UserNotifications.h>
#include <string.h>

static char *focusbar_copy_error_message(NSString *message) {
    const char *utf8 = [message UTF8String];
    if (utf8 == NULL) {
        utf8 = "unknown notification error";
    }

    return strdup(utf8);
}

char *focusbar_notify(const char *title, const char *message) {
    @autoreleasepool {
        NSString *nsTitle = [NSString stringWithUTF8String:title ? title : ""];
        NSString *nsMessage = [NSString stringWithUTF8String:message ? message : ""];

        dispatch_semaphore_t permissionSemaphore = dispatch_semaphore_create(0);
        __block BOOL granted = NO;
        __block NSError *permissionError = nil;

        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        [center requestAuthorizationWithOptions:(UNAuthorizationOptionAlert | UNAuthorizationOptionBadge | UNAuthorizationOptionSound)
                              completionHandler:^(BOOL didGrant, NSError * _Nullable error) {
            granted = didGrant;
            permissionError = error;
            dispatch_semaphore_signal(permissionSemaphore);
        }];

        dispatch_semaphore_wait(permissionSemaphore, DISPATCH_TIME_FOREVER);

        if (permissionError != nil) {
            return focusbar_copy_error_message([NSString stringWithFormat:@"authorization failed: %@", permissionError.localizedDescription]);
        }

        if (!granted) {
            return focusbar_copy_error_message(@"notification permission denied");
        }

        UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
        content.title = nsTitle;
        content.body = nsMessage;

        UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:[[NSUUID UUID] UUIDString]
                                                                              content:content
                                                                              trigger:nil];

        dispatch_semaphore_t deliverySemaphore = dispatch_semaphore_create(0);
        __block NSError *deliveryError = nil;

        [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
            deliveryError = error;
            dispatch_semaphore_signal(deliverySemaphore);
        }];

        dispatch_semaphore_wait(deliverySemaphore, DISPATCH_TIME_FOREVER);

        if (deliveryError != nil) {
            return focusbar_copy_error_message([NSString stringWithFormat:@"delivery failed: %@", deliveryError.localizedDescription]);
        }

        return NULL;
    }
}
